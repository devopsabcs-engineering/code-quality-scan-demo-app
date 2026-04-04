// =============================================================================
// validators.ts — INTENTIONAL CODE DUPLICATION
// Contains validation functions that are nearly identical to formatters.ts.
// The duplicated blocks will be detected by jscpd.
// DO NOT refactor — these violations exist for workshop demonstration.
// =============================================================================

// DUPLICATION BLOCK 1: Email validation (also in formatters.ts)
export function validateEmail(email: string): { valid: boolean; error: string | null } {
  if (!email || email.trim().length === 0) {
    return { valid: false, error: "Email is required" };
  }
  const trimmedEmail = email.trim().toLowerCase();
  if (trimmedEmail.length > 254) {
    return { valid: false, error: "Email exceeds maximum length of 254 characters" };
  }
  const atIndex = trimmedEmail.indexOf("@");
  if (atIndex === -1) {
    return { valid: false, error: "Email must contain an @ symbol" };
  }
  const localPart = trimmedEmail.substring(0, atIndex);
  const domainPart = trimmedEmail.substring(atIndex + 1);
  if (localPart.length === 0) {
    return { valid: false, error: "Email local part cannot be empty" };
  }
  if (localPart.length > 64) {
    return { valid: false, error: "Email local part exceeds 64 characters" };
  }
  if (domainPart.length === 0) {
    return { valid: false, error: "Email domain cannot be empty" };
  }
  if (domainPart.indexOf(".") === -1) {
    return { valid: false, error: "Email domain must contain a dot" };
  }
  const domainParts = domainPart.split(".");
  for (let i = 0; i < domainParts.length; i++) {
    if (domainParts[i].length === 0) {
      return { valid: false, error: "Email domain has empty label" };
    }
    if (domainParts[i].length > 63) {
      return { valid: false, error: "Email domain label exceeds 63 characters" };
    }
  }
  const lastLabel = domainParts[domainParts.length - 1];
  if (lastLabel.length < 2) {
    return { valid: false, error: "Email TLD must be at least 2 characters" };
  }
  return { valid: true, error: null };
}

// DUPLICATION BLOCK 2: Phone validation (also in formatters.ts)
export function validatePhone(phone: string): { valid: boolean; formatted: string | null; error: string | null } {
  if (!phone || phone.trim().length === 0) {
    return { valid: false, formatted: null, error: "Phone number is required" };
  }
  let cleaned = phone.replace(/[\s\-\(\)\.]/g, "");
  if (cleaned.startsWith("+")) {
    cleaned = cleaned.substring(1);
  }
  if (cleaned.startsWith("1") && cleaned.length === 11) {
    cleaned = cleaned.substring(1);
  }
  if (!/^\d+$/.test(cleaned)) {
    return { valid: false, formatted: null, error: "Phone number contains invalid characters" };
  }
  if (cleaned.length < 7) {
    return { valid: false, formatted: null, error: "Phone number is too short (minimum 7 digits)" };
  }
  if (cleaned.length > 15) {
    return { valid: false, formatted: null, error: "Phone number is too long (maximum 15 digits)" };
  }
  if (cleaned.length === 10) {
    const areaCode = cleaned.substring(0, 3);
    const exchange = cleaned.substring(3, 6);
    const subscriber = cleaned.substring(6, 10);
    const formatted = `(${areaCode}) ${exchange}-${subscriber}`;
    return { valid: true, formatted, error: null };
  }
  return { valid: true, formatted: cleaned, error: null };
}

// DUPLICATION BLOCK 3: Date validation (also in formatters.ts)
export function validateDate(dateStr: string): { valid: boolean; parsed: Date | null; error: string | null } {
  if (!dateStr || dateStr.trim().length === 0) {
    return { valid: false, parsed: null, error: "Date is required" };
  }
  const trimmedDate = dateStr.trim();
  let year: number;
  let month: number;
  let day: number;
  if (/^\d{4}-\d{2}-\d{2}$/.test(trimmedDate)) {
    const parts = trimmedDate.split("-");
    year = parseInt(parts[0], 10);
    month = parseInt(parts[1], 10);
    day = parseInt(parts[2], 10);
  } else if (/^\d{2}\/\d{2}\/\d{4}$/.test(trimmedDate)) {
    const parts = trimmedDate.split("/");
    month = parseInt(parts[0], 10);
    day = parseInt(parts[1], 10);
    year = parseInt(parts[2], 10);
  } else if (/^\d{2}\.\d{2}\.\d{4}$/.test(trimmedDate)) {
    const parts = trimmedDate.split(".");
    day = parseInt(parts[0], 10);
    month = parseInt(parts[1], 10);
    year = parseInt(parts[2], 10);
  } else {
    return { valid: false, parsed: null, error: "Unrecognized date format. Use YYYY-MM-DD, MM/DD/YYYY, or DD.MM.YYYY" };
  }
  if (month < 1 || month > 12) {
    return { valid: false, parsed: null, error: "Month must be between 1 and 12" };
  }
  const daysInMonth = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];
  if ((year % 4 === 0 && year % 100 !== 0) || year % 400 === 0) {
    daysInMonth[1] = 29;
  }
  if (day < 1 || day > daysInMonth[month - 1]) {
    return { valid: false, parsed: null, error: `Day must be between 1 and ${daysInMonth[month - 1]} for month ${month}` };
  }
  if (year < 1900 || year > 2100) {
    return { valid: false, parsed: null, error: "Year must be between 1900 and 2100" };
  }
  const parsed = new Date(year, month - 1, day);
  return { valid: true, parsed, error: null };
}

// Postal code validation
export function validatePostalCode(code: string, country: string): { valid: boolean; error: string | null } {
  if (!code || code.trim().length === 0) {
    return { valid: false, error: "Postal code is required" };
  }
  const trimmedCode = code.trim().toUpperCase();
  if (country === "US") {
    if (/^\d{5}$/.test(trimmedCode) || /^\d{5}-\d{4}$/.test(trimmedCode)) {
      return { valid: true, error: null };
    }
    return { valid: false, error: "US ZIP code must be 5 digits or 5+4 format (e.g., 12345 or 12345-6789)" };
  } else if (country === "CA") {
    if (/^[A-Z]\d[A-Z]\s?\d[A-Z]\d$/.test(trimmedCode)) {
      return { valid: true, error: null };
    }
    return { valid: false, error: "Canadian postal code must be in A1A 1A1 format" };
  } else if (country === "UK") {
    if (/^[A-Z]{1,2}\d[A-Z\d]?\s?\d[A-Z]{2}$/.test(trimmedCode)) {
      return { valid: true, error: null };
    }
    return { valid: false, error: "UK postcode format is invalid" };
  }
  return { valid: true, error: null };
}

// Credit card number validation (Luhn algorithm)
export function validateCreditCard(cardNumber: string): { valid: boolean; cardType: string | null; error: string | null } {
  if (!cardNumber || cardNumber.trim().length === 0) {
    return { valid: false, cardType: null, error: "Card number is required" };
  }
  const cleaned = cardNumber.replace(/[\s\-]/g, "");
  if (!/^\d+$/.test(cleaned)) {
    return { valid: false, cardType: null, error: "Card number must contain only digits" };
  }
  if (cleaned.length < 13 || cleaned.length > 19) {
    return { valid: false, cardType: null, error: "Card number must be between 13 and 19 digits" };
  }

  // Luhn algorithm check
  let sum = 0;
  let alternate = false;
  for (let i = cleaned.length - 1; i >= 0; i--) {
    let digit = parseInt(cleaned[i], 10);
    if (alternate) {
      digit *= 2;
      if (digit > 9) {
        digit -= 9;
      }
    }
    sum += digit;
    alternate = !alternate;
  }
  if (sum % 10 !== 0) {
    return { valid: false, cardType: null, error: "Card number failed Luhn check" };
  }

  // Card type detection
  let cardType: string | null = null;
  if (/^4/.test(cleaned)) {
    cardType = "visa";
  } else if (/^5[1-5]/.test(cleaned)) {
    cardType = "mastercard";
  } else if (/^3[47]/.test(cleaned)) {
    cardType = "amex";
  } else if (/^6(?:011|5)/.test(cleaned)) {
    cardType = "discover";
  }

  return { valid: true, cardType, error: null };
}
