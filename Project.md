# Local two-peer VPN: synthetic receive boundary

Learner-written Go IPv4 structural parsing followed by authenticated source authorization and recording-sink tests, before QUIC/TUN integration. Current exact contract covers milestone 1 only; agree subsequent policy and handler signatures before implementation.

Language: go; interface: functions.

## Milestone 1 specification: parse one IPv4 packet safely

### Deliverable and interface

Create a Go module named `localvpn` in this folder. The root package must be named `vpn`, because the supplied harness imports `localvpn`. Provide this public interface (declarations below are documentation; create the source yourself):

```go
type IPv4Info struct {
    HeaderLength int
    TotalLength  int
    Source       [4]byte
    Destination  [4]byte
}

func ParseIPv4(packet []byte) (IPv4Info, error)
```

The input is the raw bytes of exactly one complete IPv4 packet, beginning at its IPv4 header. It has no Ethernet, TUN, QUIC or other prefix. Available input is `len(packet)` bytes; spare slice capacity is not input. The function returns structural header information, not payload bytes or an authorization decision.

### Packet layout and field meanings

Offsets start at zero. Ranges below are half-open: `12:16` means offsets 12, 13, 14 and 15.

| Field | Location | Meaning |
|---|---|---|
| Version | High four bits of byte 0 | Must equal 4 |
| IHL | Low four bits of byte 0 | Header length in 4-byte words; valid values for this task are 5 through 15 |
| Total Length | Bytes `2:4` | Unsigned 16-bit big-endian integer, in bytes, including header and payload |
| Source | Bytes `12:16` | Four address octets, returned in the same order |
| Destination | Bytes `16:20` | Four address octets, returned in the same order |

`HeaderLength` is IHL converted to bytes, not the encoded word count. For example IHL 6 means a 24-byte header. A Total Length byte pair of `01 00` means 256. Options, when present, occupy the header area beyond the first 20 bytes; this task treats their contents as opaque. Payload starts after the full declared header.

### Required behavior

Reject by returning a non-nil error, without panicking, if any condition holds:

- Input length is below 20, including nil and empty input.
- Version is not 4.
- IHL is below 5.
- The declared header length exceeds available input.
- Total Length is less than the declared header length.
- Total Length differs from available input, whether the input is truncated or contains trailing bytes.

Otherwise accept: return a nil error and the exact header length in bytes, total length in bytes, source octets and destination octets. A header-only packet is allowed. Accept all inputs satisfying these structural rules; do not add checks for other header fields. Input length above 65535 necessarily cannot match the 16-bit Total Length and is rejected.

On rejection the returned `IPv4Info` fields are ignored. Error text, concrete error type and order of checks are your choice. Any panic fails the contract.

### Reproducible examples

For each row, start with the indicated number of zero bytes, set byte 0 and bytes `2:4` as shown, set source bytes `12:16` to `10,80,0,2` and destination bytes `16:20` to `10,80,0,1` wherever those offsets exist. Leave other bytes zero. Hex values are prefixed with `0x`.

| Input length | Byte 0 | Bytes `2:4` | Expected result |
|---|---|---|---|
| 20 | `0x45` | `0x00 0x14` | nil error; lengths 20 and 20; addresses as set above |
| 40 | `0x46` | `0x00 0x28` | nil error; lengths 24 and 40; addresses as set above |
| 256 | `0x45` | `0x01 0x00` | nil error; lengths 20 and 256; addresses as set above |
| 19 (capacity 64) | `0x45` | `0x00 0x13` | non-nil error; capacity does not supply a missing byte |
| 20 | `0x65` | `0x00 0x14` | non-nil error; wrong version |
| 20 | `0x44` | `0x00 0x14` | non-nil error; IHL below minimum |
| 20 | `0x46` | `0x00 0x14` | non-nil error; declared header exceeds input |
| 40 | `0x46` | `0x00 0x14` | non-nil error; total length below header length |
| 40 | `0x45` | `0x00 0x29` | non-nil error; truncated packet |
| 40 | `0x45` | `0x00 0x27` | non-nil error; trailing data |

A nil slice also returns a non-nil error. The saved `cases.json` contains additional concrete address values used by the harness.

### Scope and learner choices

No checksum verification, option interpretation, fragmentation policy, source authorization, QUIC, TUN, network access or elevated permissions is required. The synthetic examples intentionally leave unrelated fields zero. This milestone is not a complete IPv4 security validator.

Choose your source filenames, internal organization, standard-library helpers, error wording and implementation. Create the module, types and function yourself. No third-party dependencies are required. The Go toolchain must support the syntax used by the supplied harness (Go 1.18 or later). No `main` function is needed for the root library package; `harness/main.go` is the supplied test entry point.

### Run and completion criteria

After creating your module and source, run from `learn-vpn/project/local-vpn/`:

```sh
go run ./harness
```

Or from the topic folder: `cd project/local-vpn` first. Before your module or required declarations exist, module/import/undefined-symbol errors are expected. The assistant must not fill those in for you.

The harness prints one JSON report on stdout and case diagnostics on stderr. It currently selects 13 development cases. On success, all selected cases have `passed: true` and the command exits zero. Acceptance cases compare all returned fields exactly; rejection cases require a non-nil error. Panics fail. Do not change tests or expected values to accommodate a failing implementation.

Completion requires passing the development cases, then independent teacher review with fresh inputs within this same specification, followed by your explanation or trace of the important behavior. A development pass alone does not mark the milestone complete. If you choose source filenames other than `ipv4.go`, update `sourcePaths` in `project.json` before review so all implementation files are included.

## Future milestone outline (not ready for implementation)

### Authorize sources before packet injection

After structural milestone, agree authenticated-peer policy representation, handler/sink interface and exact acceptance/rejection behavior. Implement source authorization plus recording-sink tests proving unauthorized packets cause zero writes. Do not begin until contract and fixtures agreed.

Next action: Agree source policy and receive-handler interface after milestone 1 passes review and understanding check.

## Test command

Run from this folder. See the build-lesson reference for the harness report contract.

```json
["go","run","./harness"]
```

## Learner-owned setup

From the topic folder, your workspace is `project/local-vpn/`. Only the harness, cases and brief/metadata are supplied. You choose your source filenames and create the module and code yourself. The documented return type has `HeaderLength int`, `TotalLength int`, `Source [4]byte`, and `Destination [4]byte`; the harness calls `ParseIPv4(packet []byte) (IPv4Info, error)` from package `vpn`. If you choose different source filenames, update `sourcePaths` before review.

The original untouched assistant stub and manifest are recoverable in `Learning/.pi/backups/2026-10-23-vpn-scaffolding/`. `projects/` is a compatibility alias for Pi and old lesson links. No correctness or understanding check has been completed.
