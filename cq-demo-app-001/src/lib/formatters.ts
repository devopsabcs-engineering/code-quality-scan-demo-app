// =============================================================================
// formatters.ts — INTENTIONAL CODE DUPLICATION
// Contains formatting functions that are nearly identical to validators.ts.
// The duplicated blocks will be detected by jscpd.
// DO NOT refactor — these violations exist for workshop demonstration.
// =============================================================================

// DUPLICATION BLOCK 1: Email formatting — nearly identical to validateEmail in validators.ts
export function formatAndValidateEmail(email: string): { valid: boolean; formatted: string | null; error: string | null } {
  if (!email || email.trim().length === 0) {
    return { valid: false, formatted: null, error: "Email is required" };
  }
  const trimmedEmail = email.trim().toLowerCase();
  if (trimmedEmail.length > 254) {
    return { valid: false, formatted: null, error: "Email exceeds maximum length of 254 characters" };
  }
  const atIndex = trimmedEmail.indexOf("@");
  if (atIndex === -1) {
    return { valid: false, formatted: null, error: "Email must contain an @ symbol" };
  }
  const localPart = trimmedEmail.substring(0, atIndex);
  const domainPart = trimmedEmail.substring(atIndex + 1);
  if (localPart.length === 0) {
    return { valid: false, formatted: null, error: "Email local part cannot be empty" };
  }
  if (localPart.length > 64) {
    return { valid: false, formatted: null, error: "Email local part exceeds 64 characters" };
  }
  if (domainPart.length === 0) {
    return { valid: false, formatted: null, error: "Email domain cannot be empty" };
  }
  if (domainPart.indexOf(".") === -1) {
    return { valid: false, formatted: null, error: "Email domain must contain a dot" };
  }
  const domainParts = domainPart.split(".");
  for (let i = 0; i < domainParts.length; i++) {
    if (domainParts[i].length === 0) {
      return { valid: false, formatted: null, error: "Email domain has empty label" };
    }
    if (domainParts[i].length > 63) {
      return { valid: false, formatted: null, error: "Email domain label exceeds 63 characters" };
    }
  }
  const lastLabel = domainParts[domainParts.length - 1];
  if (lastLabel.length < 2) {
    return { valid: false, formatted: null, error: "Email TLD must be at least 2 characters" };
  }
  return { valid: true, formatted: trimmedEmail, error: null };
}

// DUPLICATION BLOCK 2: Phone formatting — nearly identical to validatePhone in validators.ts
export function formatPhone(phone: string): { valid: boolean; formatted: string | null; error: string | null } {
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

// DUPLICATION BLOCK 3: Date formatting — nearly identical to validateDate in validators.ts
export function formatDate(dateStr: string, outputFormat: string = "iso"): { valid: boolean; formatted: string | null; error: string | null } {
  if (!dateStr || dateStr.trim().length === 0) {
    return { valid: false, formatted: null, error: "Date is required" };
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
    return { valid: false, formatted: null, error: "Unrecognized date format. Use YYYY-MM-DD, MM/DD/YYYY, or DD.MM.YYYY" };
  }
  if (month < 1 || month > 12) {
    return { valid: false, formatted: null, error: "Month must be between 1 and 12" };
  }
  const daysInMonth = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];
  if ((year % 4 === 0 && year % 100 !== 0) || year % 400 === 0) {
    daysInMonth[1] = 29;
  }
  if (day < 1 || day > daysInMonth[month - 1]) {
    return { valid: false, formatted: null, error: `Day must be between 1 and ${daysInMonth[month - 1]} for month ${month}` };
  }
  if (year < 1900 || year > 2100) {
    return { valid: false, formatted: null, error: "Year must be between 1900 and 2100" };
  }

  const monthStr = month.toString().padStart(2, "0");
  const dayStr = day.toString().padStart(2, "0");

  let formatted: string;
  if (outputFormat === "us") {
    formatted = `${monthStr}/${dayStr}/${year}`;
  } else if (outputFormat === "eu") {
    formatted = `${dayStr}.${monthStr}.${year}`;
  } else {
    formatted = `${year}-${monthStr}-${dayStr}`;
  }

  return { valid: true, formatted, error: null };
}

// Currency formatting
export function formatCurrency(amount: number, currency: string = "USD"): string {
  const currencySymbols: Record<string, string> = {
    USD: "$",
    EUR: "€",
    GBP: "£",
    JPY: "¥",
    CAD: "C$",
    AUD: "A$",
  };
  const symbol = currencySymbols[currency] || currency + " ";
  const decimals = currency === "JPY" ? 0 : 2;
  const formatted = Math.abs(amount).toFixed(decimals).replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  if (amount < 0) {
    return `-${symbol}${formatted}`;
  }
  return `${symbol}${formatted}`;
}

// Address formatting
export function formatAddress(
  street: string,
  city: string,
  state: string,
  postalCode: string,
  country: string
): string {
  const parts: string[] = [];
  if (street && street.trim()) {
    parts.push(street.trim());
  }
  if (city && city.trim()) {
    let cityLine = city.trim();
    if (state && state.trim()) {
      cityLine += `, ${state.trim()}`;
    }
    if (postalCode && postalCode.trim()) {
      cityLine += ` ${postalCode.trim()}`;
    }
    parts.push(cityLine);
  }
  if (country && country.trim()) {
    parts.push(country.trim());
  }
  return parts.join("\n");
}

// Percentage formatting
export function formatPercentage(value: number, decimals: number = 1): string {
  if (isNaN(value) || !isFinite(value)) {
    return "N/A";
  }
  return `${value.toFixed(decimals)}%`;
}
