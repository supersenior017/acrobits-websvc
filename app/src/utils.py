def to_lower_camel_case(string: str) -> str:
    first, *rest = string.split('_')
    return f'{first}{"".join(map(str.capitalize, rest))}'
