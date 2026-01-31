#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Unit tests for GTP5G Netlink Decoder

Run with: python3 -m pytest test_netlink_decoder.py -v
Or: python3 -m unittest test_netlink_decoder -v
"""

import unittest
import struct
import io
import sys
from unittest.mock import patch

# Import functions from the decoder module
from netlink_decoder import (
    decode_value,
    parse_attributes,
    parse_nlm_flags,
    parse_nlmsg_type,
    format_attrs,
    process_line,
    get_gtp5g_family_id,
    GTP5G_PDR_ATTRS,
    GTP5G_FAR_ATTRS,
    GTP5G_QER_ATTRS,
    GTP5G_URR_ATTRS,
    GTP5G_BAR_ATTRS,
    GTP5G_PDI_ATTRS,
    GTP5G_F_TEID_ATTRS,
    GTP5G_COMMON_ATTRS,
    NLM_FLAGS,
)


class TestDecodeValue(unittest.TestCase):
    """Tests for decode_value function."""

    def test_decode_ipv4_address(self):
        """Test decoding IPv4 addresses."""
        # 192.168.1.1
        data = bytes([192, 168, 1, 1])
        result = decode_value("GTP5G_PDR_ROLE_ADDR_IPV4", data)
        self.assertEqual(result, "192.168.1.1")
        
        # 10.0.0.1
        data = bytes([10, 0, 0, 1])
        result = decode_value("GTP5G_F_TEID_GTPU_ADDR_IPV4", data)
        self.assertEqual(result, "10.0.0.1")
        
        # 0.0.0.0
        data = bytes([0, 0, 0, 0])
        result = decode_value("GTP5G_PDI_UE_ADDR_IPV4", data)
        self.assertEqual(result, "0.0.0.0")

    def test_decode_u8_values(self):
        """Test decoding U8 values."""
        # QER Gate status
        data = struct.pack("=B", 3)
        result = decode_value("GTP5G_QER_GATE", data)
        self.assertEqual(result, 3)
        
        # Outer header removal
        data = struct.pack("=B", 0)
        result = decode_value("GTP5G_OUTER_HEADER_REMOVAL", data)
        self.assertEqual(result, 0)
        
        # QFI value
        data = struct.pack("=B", 9)
        result = decode_value("GTP5G_QER_QFI", data)
        self.assertEqual(result, 9)

    def test_decode_u16_values(self):
        """Test decoding U16 values."""
        # PDR ID
        data = struct.pack("=H", 1234)
        result = decode_value("GTP5G_PDR_ID", data)
        self.assertEqual(result, 1234)
        
        # FAR Apply Action
        data = struct.pack("=H", 2)
        result = decode_value("GTP5G_FAR_APPLY_ACTION", data)
        self.assertEqual(result, 2)
        
        # Port number
        data = struct.pack("=H", 2152)
        result = decode_value("GTP5G_OUTER_HEADER_CREATION_PORT", data)
        self.assertEqual(result, 2152)

    def test_decode_u32_values(self):
        """Test decoding U32 values."""
        # FAR ID
        data = struct.pack("=I", 1)
        result = decode_value("GTP5G_FAR_ID", data)
        self.assertEqual(result, 1)
        
        # TEID
        data = struct.pack("=I", 0x12345678)
        result = decode_value("GTP5G_F_TEID_I_TEID", data)
        self.assertEqual(result, 0x12345678)
        
        # Link interface index
        data = struct.pack("=I", 5)
        result = decode_value("GTP5G_LINK", data)
        self.assertEqual(result, 5)
        
        # Precedence
        data = struct.pack("=I", 100)
        result = decode_value("GTP5G_PDR_PRECEDENCE", data)
        self.assertEqual(result, 100)

    def test_decode_u64_values(self):
        """Test decoding U64 values (SEID, timestamps)."""
        # SEID
        data = struct.pack("=Q", 0x123456789ABCDEF0)
        result = decode_value("GTP5G_PDR_SEID", data)
        self.assertEqual(result, 0x123456789ABCDEF0)
        
        # Volume threshold
        data = struct.pack("=Q", 1000000)
        result = decode_value("GTP5G_URR_VOLUME_THRESHOLD_TOVOL", data)
        self.assertEqual(result, 1000000)

    def test_decode_string_values(self):
        """Test decoding string values."""
        # Unix socket path
        data = b"/tmp/upf.sock\x00"
        result = decode_value("GTP5G_PDR_UNIX_SOCKET_PATH", data)
        self.assertEqual(result, "/tmp/upf.sock")
        
        # Forwarding policy
        data = b"policy1\x00"
        result = decode_value("GTP5G_FORWARDING_PARAMETER_FORWARDING_POLICY", data)
        self.assertEqual(result, "policy1")

    def test_decode_port_ranges(self):
        """Test decoding port range values."""
        # Single port (80-80)
        data = struct.pack("=I", (80 << 16) | 80)
        result = decode_value("GTP5G_FLOW_DESCRIPTION_SRC_PORT", data)
        self.assertEqual(result, "80")
        
        # Port range (1024-2048)
        data = struct.pack("=I", (2048 << 16) | 1024)
        result = decode_value("GTP5G_FLOW_DESCRIPTION_DEST_PORT", data)
        self.assertEqual(result, "1024-2048")
        
        # Empty port
        result = decode_value("GTP5G_FLOW_DESCRIPTION_SRC_PORT", b"")
        self.assertEqual(result, "(none)")

    def test_decode_mask_values(self):
        """Test decoding network mask values (big-endian)."""
        # 255.255.255.0
        data = struct.pack(">I", 0xFFFFFF00)
        result = decode_value("GTP5G_FLOW_DESCRIPTION_SRC_MASK", data)
        self.assertEqual(result, "255.255.255.0")
        
        # 255.255.0.0
        data = struct.pack(">I", 0xFFFF0000)
        result = decode_value("GTP5G_FLOW_DESCRIPTION_DEST_MASK", data)
        self.assertEqual(result, "255.255.0.0")

    def test_decode_unknown_attribute(self):
        """Test decoding unknown attribute returns hex string."""
        data = bytes([0xDE, 0xAD, 0xBE, 0xEF])
        result = decode_value("UNKNOWN_ATTR", data)
        self.assertEqual(result, "0xdeadbeef")

    def test_decode_empty_data(self):
        """Test decoding empty data."""
        result = decode_value("GTP5G_LINK", b"")
        self.assertEqual(result, "(empty)")

    def test_decode_insufficient_data(self):
        """Test decoding with insufficient data length."""
        # U32 requires 4 bytes, only provide 2
        data = bytes([0x01, 0x02])
        result = decode_value("GTP5G_FAR_ID", data)
        # Should return hex fallback
        self.assertEqual(result, "0x0102")


class TestParseNlmFlags(unittest.TestCase):
    """Tests for parse_nlm_flags function."""

    def test_parse_numeric_flags(self):
        """Test parsing numeric flag values."""
        self.assertEqual(parse_nlm_flags("5"), 5)
        self.assertEqual(parse_nlm_flags("0x05"), 5)
        self.assertEqual(parse_nlm_flags("0x300"), 0x300)

    def test_parse_single_symbolic_flag(self):
        """Test parsing single symbolic flags."""
        self.assertEqual(parse_nlm_flags("NLM_F_REQUEST"), 0x01)
        self.assertEqual(parse_nlm_flags("NLM_F_ACK"), 0x04)
        self.assertEqual(parse_nlm_flags("NLM_F_MULTI"), 0x02)

    def test_parse_combined_symbolic_flags(self):
        """Test parsing combined symbolic flags."""
        result = parse_nlm_flags("NLM_F_REQUEST|NLM_F_ACK")
        self.assertEqual(result, 0x05)
        
        result = parse_nlm_flags("NLM_F_REQUEST|NLM_F_MULTI|NLM_F_ACK")
        self.assertEqual(result, 0x07)

    def test_parse_mixed_flags(self):
        """Test parsing mixed symbolic and numeric flags."""
        result = parse_nlm_flags("NLM_F_REQUEST|0x200")
        self.assertEqual(result, 0x201)

    def test_parse_flags_with_whitespace(self):
        """Test parsing flags with extra whitespace."""
        result = parse_nlm_flags("  NLM_F_REQUEST | NLM_F_ACK  ")
        self.assertEqual(result, 0x05)


class TestParseNlmsgType(unittest.TestCase):
    """Tests for parse_nlmsg_type function."""

    def test_parse_numeric_type(self):
        """Test parsing numeric message types."""
        self.assertEqual(parse_nlmsg_type("31"), 31)
        self.assertEqual(parse_nlmsg_type("0x1f"), 31)

    def test_parse_type_with_comment(self):
        """Test parsing type with NLMSG_??? comment."""
        result = parse_nlmsg_type("0x1f /* NLMSG_??? */")
        self.assertEqual(result, 31)
        
        result = parse_nlmsg_type("31 /* GENERIC_FAMILY_??? */")
        self.assertEqual(result, 31)

    def test_parse_symbolic_type(self):
        """Test parsing symbolic types."""
        self.assertEqual(parse_nlmsg_type("NLMSG_OVERRUN"), 4)

    def test_parse_unknown_type(self):
        """Test parsing unknown type returns 0."""
        self.assertEqual(parse_nlmsg_type("UNKNOWN_TYPE"), 0)


class TestParseAttributes(unittest.TestCase):
    """Tests for parse_attributes function."""

    def test_parse_empty_data(self):
        """Test parsing empty attribute data."""
        result = parse_attributes(b"", GTP5G_COMMON_ATTRS)
        self.assertEqual(result, {})

    def test_parse_single_u32_attribute(self):
        """Test parsing single U32 attribute."""
        # NLA: len=8, type=1 (GTP5G_LINK), value=5
        data = struct.pack("=HH", 8, 1) + struct.pack("=I", 5)
        result = parse_attributes(data, GTP5G_COMMON_ATTRS)
        self.assertEqual(result.get("GTP5G_LINK"), 5)

    def test_parse_single_u16_attribute(self):
        """Test parsing single U16 attribute (PDR_ID)."""
        # NLA: len=6, type=3 (GTP5G_PDR_ID), value=1234 + 2 bytes padding
        data = struct.pack("=HH", 6, 3) + struct.pack("=H", 1234) + b"\x00\x00"
        result = parse_attributes(data, GTP5G_PDR_ATTRS)
        self.assertEqual(result.get("GTP5G_PDR_ID"), 1234)

    def test_parse_multiple_attributes(self):
        """Test parsing multiple attributes."""
        # Attribute 1: GTP5G_LINK (type=1), value=5
        attr1 = struct.pack("=HH", 8, 1) + struct.pack("=I", 5)
        # Attribute 2: GTP5G_NET_NS_FD (type=2), value=10
        attr2 = struct.pack("=HH", 8, 2) + struct.pack("=I", 10)
        
        data = attr1 + attr2
        result = parse_attributes(data, GTP5G_COMMON_ATTRS)
        
        self.assertEqual(result.get("GTP5G_LINK"), 5)
        self.assertEqual(result.get("GTP5G_NET_NS_FD"), 10)

    def test_parse_ipv4_attribute(self):
        """Test parsing IPv4 address attribute."""
        # NLA: len=8, type=8 (GTP5G_PDR_ROLE_ADDR_IPV4), value=192.168.1.1
        data = struct.pack("=HH", 8, 8) + bytes([192, 168, 1, 1])
        result = parse_attributes(data, GTP5G_PDR_ATTRS)
        self.assertEqual(result.get("GTP5G_PDR_ROLE_ADDR_IPV4"), "192.168.1.1")

    def test_parse_nested_attribute(self):
        """Test parsing nested attribute (PDI)."""
        # Inner attribute: GTP5G_PDI_SRC_INTF (type=4), value=1
        inner_attr = struct.pack("=HH", 5, 4) + struct.pack("=B", 1) + b"\x00\x00\x00"
        
        # Outer attribute: GTP5G_PDR_PDI (type=5) with NLA_F_NESTED flag
        outer_len = 4 + len(inner_attr)
        outer_attr = struct.pack("=HH", outer_len, 5 | 0x8000) + inner_attr
        
        result = parse_attributes(outer_attr, GTP5G_PDR_ATTRS)
        
        self.assertIn("GTP5G_PDR_PDI", result)
        self.assertIsInstance(result["GTP5G_PDR_PDI"], dict)
        self.assertEqual(result["GTP5G_PDR_PDI"].get("GTP5G_PDI_SRC_INTF"), 1)

    def test_parse_with_alignment(self):
        """Test that attribute parsing handles 4-byte alignment."""
        # 5-byte attribute (needs padding to 8)
        # len=5, type=4 (SRC_INTF), value=1
        attr1 = struct.pack("=HH", 5, 4) + struct.pack("=B", 1)
        # Padding to align to 4 bytes
        attr1_padded = attr1 + b"\x00\x00\x00"
        
        # Next attribute
        attr2 = struct.pack("=HH", 8, 1) + struct.pack("=I", 99)
        
        data = attr1_padded + attr2
        result = parse_attributes(data, GTP5G_PDI_ATTRS)
        
        self.assertEqual(result.get("GTP5G_PDI_SRC_INTF"), 1)

    def test_parse_unknown_attribute(self):
        """Test parsing unknown attribute type."""
        # NLA: len=8, type=99 (unknown), value=0xDEADBEEF
        data = struct.pack("=HH", 8, 99) + struct.pack("=I", 0xDEADBEEF)
        result = parse_attributes(data, GTP5G_COMMON_ATTRS)
        self.assertIn("UNKNOWN_ATTR_99", result)


class TestFormatAttrs(unittest.TestCase):
    """Tests for format_attrs function."""

    def test_format_empty_attrs(self):
        """Test formatting empty attributes."""
        result = format_attrs({})
        self.assertEqual(result, "  (empty)")

    def test_format_simple_attrs(self):
        """Test formatting simple attributes."""
        attrs = {"GTP5G_LINK": 5, "GTP5G_FAR_ID": 1}
        result = format_attrs(attrs)
        self.assertIn("GTP5G_LINK: 5", result)
        self.assertIn("GTP5G_FAR_ID: 1", result)

    def test_format_nested_attrs(self):
        """Test formatting nested attributes."""
        attrs = {
            "GTP5G_PDR_PDI": {
                "GTP5G_PDI_SRC_INTF": 0,
                "GTP5G_PDI_UE_ADDR_IPV4": "10.60.0.1"
            }
        }
        result = format_attrs(attrs)
        self.assertIn("GTP5G_PDR_PDI:", result)
        self.assertIn("GTP5G_PDI_SRC_INTF: 0", result)
        self.assertIn("GTP5G_PDI_UE_ADDR_IPV4: 10.60.0.1", result)


class TestProcessLine(unittest.TestCase):
    """Tests for process_line function."""

    def setUp(self):
        """Set up test fixtures."""
        self.gtp5g_family_id = 31

    def test_process_non_netlink_line(self):
        """Test that non-netlink lines are ignored."""
        # Capture stdout
        captured = io.StringIO()
        with patch('sys.stdout', captured):
            process_line("some random log line", self.gtp5g_family_id)
        
        self.assertEqual(captured.getvalue(), "")

    def test_process_unfinished_line(self):
        """Test that unfinished syscalls are ignored."""
        captured = io.StringIO()
        with patch('sys.stdout', captured):
            process_line("sendmsg(5, <unfinished ...>", self.gtp5g_family_id)
        
        self.assertEqual(captured.getvalue(), "")

    def test_process_nlmsg_error_line(self):
        """Test that NLMSG_ERROR responses are skipped."""
        line = 'recvmsg(5, {msg_iov=[{iov_base={type=NLMSG_ERROR, ...}}'
        captured = io.StringIO()
        with patch('sys.stdout', captured):
            process_line(line, self.gtp5g_family_id)
        
        self.assertEqual(captured.getvalue(), "")

    def test_process_gtp5g_del_far_message(self):
        """Test processing a GTP5G_CMD_DEL_FAR message."""
        # Real strace output format for DEL_FAR command
        # GenL header: cmd=5 (DEL_FAR), version=0, reserved=0
        # Attributes: LINK=5, FAR_ID=1
        genl_header = struct.pack("=BBH", 5, 0, 0)  # cmd=5, ver=0, reserved=0
        attr_link = struct.pack("=HH", 8, 1) + struct.pack("=I", 5)  # LINK=5
        attr_far_id = struct.pack("=HH", 8, 3) + struct.pack("=I", 1)  # FAR_ID=1
        
        payload = genl_header + attr_link + attr_far_id
        hex_payload = ''.join(f'\\x{b:02x}' for b in payload)
        
        # Real strace format: single brace with separate iov entries
        # First iov has header, second iov has GenL payload
        line = (
            f'sendmsg(5, {{msg_iov=[{{iov_base={{len=36, type=gtp5g, '
            f'flags=NLM_F_REQUEST|NLM_F_ACK, seq=1, pid=0}}, iov_len=16}}, '
            f'{{iov_base="{hex_payload}", iov_len=20}}], msg_iovlen=2}}, 0) = 36'
        )
        
        captured = io.StringIO()
        with patch('sys.stdout', captured):
            process_line(line, self.gtp5g_family_id)
        
        output = captured.getvalue()
        self.assertIn("GTP5G MESSAGE", output)
        self.assertIn("GTP5G_CMD_DEL_FAR", output)
        self.assertIn("GTP5G_LINK: 5", output)
        self.assertIn("GTP5G_FAR_ID: 1", output)

    def test_process_different_family_id_ignored(self):
        """Test that messages with different family ID are ignored."""
        # Use family ID 16 (not gtp5g's 31)
        line = 'sendmsg(5, {msg_iov=[{iov_base={len=20, type=0x10, flags=5, seq=1, pid=0}}], msg_iovlen=1}, 0) = 20'
        
        captured = io.StringIO()
        with patch('sys.stdout', captured):
            process_line(line, self.gtp5g_family_id)
        
        self.assertEqual(captured.getvalue(), "")

    def test_process_gtp5g_add_pdr_message(self):
        """Test processing a GTP5G_CMD_ADD_PDR message."""
        # GenL header: cmd=1 (ADD_PDR), version=0, reserved=0
        genl_header = struct.pack("=BBH", 1, 0, 0)
        # Attributes
        attr_link = struct.pack("=HH", 8, 1) + struct.pack("=I", 5)  # LINK=5
        attr_pdr_id = struct.pack("=HH", 6, 3) + struct.pack("=H", 100) + b"\x00\x00"  # PDR_ID=100
        attr_precedence = struct.pack("=HH", 8, 4) + struct.pack("=I", 255)  # PRECEDENCE=255
        
        payload = genl_header + attr_link + attr_pdr_id + attr_precedence
        hex_payload = ''.join(f'\\x{b:02x}' for b in payload)
        
        # Real strace format: single brace with separate iov entries
        line = (
            f'sendmsg(5, {{msg_iov=[{{iov_base={{len=48, type=gtp5g, '
            f'flags=NLM_F_REQUEST|NLM_F_ACK, seq=2, pid=0}}, iov_len=16}}, '
            f'{{iov_base="{hex_payload}", iov_len=28}}], msg_iovlen=2}}, 0) = 48'
        )
        
        captured = io.StringIO()
        with patch('sys.stdout', captured):
            process_line(line, self.gtp5g_family_id)
        
        output = captured.getvalue()
        self.assertIn("GTP5G_CMD_ADD_PDR", output)
        self.assertIn("GTP5G_PDR_ID: 100", output)
        self.assertIn("GTP5G_PDR_PRECEDENCE: 255", output)


class TestGetGtp5gFamilyId(unittest.TestCase):
    """Tests for get_gtp5g_family_id function."""

    @patch('subprocess.run')
    def test_get_family_id_success(self, mock_run):
        """Test successful family ID detection."""
        mock_run.return_value.returncode = 0
        mock_run.return_value.stdout = "Name: gtp5g\n  ID: 0x1f\n"
        
        result = get_gtp5g_family_id()
        self.assertEqual(result, 31)

    @patch('subprocess.run')
    def test_get_family_id_single_line(self, mock_run):
        """Test family ID detection with single-line format."""
        mock_run.return_value.returncode = 0
        mock_run.return_value.stdout = "Name: gtp5g ID: 0x1f"
        
        result = get_gtp5g_family_id()
        self.assertEqual(result, 31)

    @patch('subprocess.run')
    def test_get_family_id_command_fails(self, mock_run):
        """Test handling when genl command fails."""
        mock_run.return_value.returncode = 1
        mock_run.return_value.stdout = ""
        
        result = get_gtp5g_family_id()
        self.assertIsNone(result)

    @patch('subprocess.run')
    def test_get_family_id_not_found(self, mock_run):
        """Test handling when gtp5g module not loaded."""
        mock_run.return_value.returncode = 0
        mock_run.return_value.stdout = "No such Generic Netlink family"
        
        result = get_gtp5g_family_id()
        self.assertIsNone(result)

    @patch('subprocess.run', side_effect=FileNotFoundError())
    def test_get_family_id_genl_missing(self, mock_run):
        """Test handling when genl command not installed."""
        result = get_gtp5g_family_id()
        self.assertIsNone(result)


class TestIntegration(unittest.TestCase):
    """Integration tests with realistic strace output."""

    def test_full_strace_line_del_far(self):
        """Test processing complete strace line for DEL_FAR."""
        # Use actual strace output sample format
        # GenL header in separate iov: cmd=5 (DEL_FAR)
        genl_cmd = '\\x05\\x00\\x00\\x00'  # cmd=5, ver=0, reserved=0
        # Attributes: LINK and FAR_ID in another iov
        attr_data = '\\x0c\\x00\\x07\\x00\\x01\\x00\\x00\\x00\\x00\\x00\\x00\\x00'
        
        line = (
            f'[pid 24340] sendmsg(13, {{msg_name={{sa_family=AF_NETLINK, nl_pid=0, nl_groups=00000000}}, '
            f'msg_namelen=12, msg_iov=[{{iov_base={{len=48, type=gtp5g, flags=NLM_F_REQUEST|NLM_F_ACK|0x200, '
            f'seq=65, pid=0}}, iov_len=16}}, {{iov_base="{genl_cmd}", iov_len=4}}, '
            f'{{iov_base={{len=65544, type=0x7 /* NLMSG_??? */, flags=0, seq=196616, pid=8}}, iov_len=16}}, '
            f'{{iov_base="{attr_data}", iov_len=12}}], msg_iovlen=4, msg_controllen=0, msg_flags=0}}, 0) = 48'
        )
        
        captured = io.StringIO()
        with patch('sys.stdout', captured):
            process_line(line, 31)
        
        output = captured.getvalue()
        self.assertIn("GTP5G MESSAGE", output)
        self.assertIn("GTP5G_CMD_DEL_FAR", output)

    def test_strace_line_with_single_brace_format(self):
        """Test strace output with single brace iov_base format."""
        line = (
            'sendmsg(5, {msg_iov=[{iov_base={len=36, type=0x1f, flags=5, seq=1, pid=0}}], '
            'msg_iovlen=1}, 0) = 36'
        )
        
        captured = io.StringIO()
        with patch('sys.stdout', captured):
            # This should match but have no payload to decode
            process_line(line, 31)
        
        # The message is matched but no hex payload, so minimal output
        # This tests that the regex pattern works


class TestEdgeCases(unittest.TestCase):
    """Tests for edge cases and error handling."""

    def test_malformed_hex_data(self):
        """Test handling of malformed hex data in iov_base."""
        line = 'sendmsg(5, {msg_iov=[{iov_base="\\xZZ\\xYY"}]}, 0) = 4'
        
        # Should not crash
        captured = io.StringIO()
        with patch('sys.stdout', captured):
            process_line(line, 31)

    def test_very_short_payload(self):
        """Test handling of payload shorter than GenL header."""
        line = 'sendmsg(5, {msg_iov=[{iov_base={len=4, type=0x1f, flags=5, seq=1, pid=0}, "\\x01\\x02"}]}, 0)'
        
        captured = io.StringIO()
        with patch('sys.stdout', captured):
            process_line(line, 31)

    def test_zero_length_nla(self):
        """Test that zero-length NLA is handled."""
        # GenL header + zero-length attribute
        payload = struct.pack("=BBH", 5, 0, 0) + struct.pack("=HH", 0, 1)
        
        # parse_attributes should skip zero-length
        from netlink_decoder import parse_attributes
        result = parse_attributes(payload[4:], GTP5G_COMMON_ATTRS)
        # Should not crash, may have empty result


if __name__ == "__main__":
    unittest.main(verbosity=2)
