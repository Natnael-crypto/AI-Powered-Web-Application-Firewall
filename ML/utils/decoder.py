import re
import urllib.parse
import html
import base64

def recursive_decode(input_string: str, max_depth: int) -> str:
    current = input_string

    for i in range(max_depth):
        decoded = current

        decoded = url_decode(decoded)
        decoded = html_decode(decoded)
        decoded = unicode_decode(decoded)
        decoded = hex_decode(decoded)
        decoded = base64_decode(decoded)

        decoded = normalize_payload(decoded)

        if decoded == current:
            break
        current = decoded
    return current

def url_decode(s: str) -> str:
    try:
        decoded = urllib.parse.unquote(s)
        return decoded
    except:
        return s

def html_decode(s: str) -> str:
    return html.unescape(s)

def unicode_decode(s: str) -> str:
    def replace_func(match):
        hex_val = match.group(2)
        return chr(int(hex_val, 16))

    re_unicode = re.compile(r'(\\u|%u)([0-9a-fA-F]{4})')
    return re_unicode.sub(replace_func, s)

def hex_decode(s: str) -> str:
    # Decode \xHH sequences
    def replace_func_x(match):
        hex_byte = match.group(1)
        return chr(int(hex_byte, 16))

    re_hex_x = re.compile(r'\\x([0-9a-fA-F]{2})')
    s = re_hex_x.sub(replace_func_x, s)

    # Decode a full hex string if it's a valid hex representation
    if len(s) % 2 == 0 and is_hex_string(s):
        try:
            return bytes.fromhex(s).decode('utf-8', errors='ignore')
        except:
            pass # Fall through to return s if decoding fails
    return s

def base64_decode(s: str) -> str:
    s = s.strip()
    if len(s) < 8 or len(s) % 4 != 0:
        return s
    try:
        decoded_bytes = base64.b64decode(s)
        if is_mostly_printable(decoded_bytes):
            return decoded_bytes.decode('utf-8', errors='ignore')
        return s
    except:
        return s

def is_mostly_printable(data: bytes) -> bool:
    printable = 0
    for b in data:
        if b in (9, 10, 13) or (32 <= b <= 126):
            printable += 1
    return float(printable) / len(data) > 0.8 if data else False

def is_hex_string(s: str) -> bool:
    re_hex_only = re.compile(r'^[0-9a-fA-F]+$')
    return bool(re_hex_only.match(s))

def normalize_payload(s: str) -> str:
    s = s.lower()
    s = s.replace('\n', ' ')
    s = s.replace('\r', ' ')
    s = re.sub(r'\s+', ' ', s)
    return s