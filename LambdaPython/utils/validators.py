def validate_string(value):
    if not isinstance(value, str):
        raise ValueError("must be a string")

    value = value.strip()

    if not value:
        raise ValueError("must not be empty")

    return value


def validate_emp_id(value):
    value = validate_string(value)

    if not value.startswith("EMP"):
        raise ValueError("must start with EMP")

    if not value[3:].isdigit():
        raise ValueError("must contain only digits after EMP")

    if len(value) < 6:
        raise ValueError("must be at least 6 characters long")

    return value


def validate_contact_no(value):
    value = validate_string(value)

    if not value.isdigit():
        raise ValueError("must contain only digits")

    if len(value) < 10:
        raise ValueError("must be at least 10 digits")

    if len(value) > 12:
        raise ValueError("must not exceed 12 digits")

    return value