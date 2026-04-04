// =============================================================================
// helpers.ts — INTENTIONAL LINT VIOLATIONS
// Contains ESLint violations: unused variables, `let` instead of `const`,
// missing return types, `any` types, and other style issues.
// DO NOT fix — these violations exist for workshop demonstration.
// =============================================================================

// VIOLATION: prefer-const — `let` used but never reassigned
let APP_NAME = "cq-demo-app-001";

// VIOLATION: prefer-const — `let` used but never reassigned
let API_VERSION = "v1";

// VIOLATION: prefer-const — `let` used but never reassigned
let MAX_PAGE_SIZE = 100;

// VIOLATION: prefer-const — `let` used but never reassigned
let DEFAULT_TIMEOUT = 30000;

// VIOLATION: @typescript-eslint/no-unused-vars — variable declared but never used
let debugMode = false;

// VIOLATION: @typescript-eslint/no-unused-vars — variable declared but never used
let internalCounter = 0;

// VIOLATION: @typescript-eslint/no-unused-vars — variable declared but never used
let temporaryBuffer: string[] = [];

// VIOLATION: @typescript-eslint/no-unused-vars — function declared but never used
function logInternal(message: string, level: string): void {
  // This function is intentionally unused
  const timestamp = new Date().toISOString();
  console.log(`[${timestamp}] [${level}] ${message}`);
}

// VIOLATION: @typescript-eslint/no-explicit-any — uses `any` type
export function parseApiResponse(response: any): any {
  // VIOLATION: prefer-const — `let` used but never reassigned
  let data = response.data;
  // VIOLATION: prefer-const
  let status = response.status;
  // VIOLATION: @typescript-eslint/no-unused-vars — variable assigned but never used
  let metadata = response.headers;

  if (status >= 200 && status < 300) {
    return { success: true, data: data };
  } else if (status >= 400 && status < 500) {
    return { success: false, error: "Client error", code: status };
  } else {
    return { success: false, error: "Server error", code: status };
  }
}

// VIOLATION: @typescript-eslint/no-explicit-any — uses `any` type
export function mergeObjects(target: any, source: any): any {
  // VIOLATION: prefer-const
  let result = { ...target };
  for (let key in source) {
    if (source.hasOwnProperty(key)) {
      if (typeof source[key] === "object" && source[key] !== null && !Array.isArray(source[key])) {
        result[key] = mergeObjects(result[key] || {}, source[key]);
      } else {
        result[key] = source[key];
      }
    }
  }
  return result;
}

// VIOLATION: @typescript-eslint/no-explicit-any — uses `any` type throughout
export function buildQueryString(params: any): string {
  // VIOLATION: prefer-const
  let parts: string[] = [];
  for (let key in params) {
    if (params[key] !== undefined && params[key] !== null) {
      // VIOLATION: prefer-const
      let value = String(params[key]);
      parts.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`);
    }
  }
  return parts.length > 0 ? `?${parts.join("&")}` : "";
}

// VIOLATION: @typescript-eslint/no-explicit-any — parameter and return type
export function flattenArray(arr: any[]): any[] {
  // VIOLATION: prefer-const
  let result: any[] = [];
  for (let i = 0; i < arr.length; i++) {
    if (Array.isArray(arr[i])) {
      // VIOLATION: prefer-const
      let nested = flattenArray(arr[i]);
      for (let j = 0; j < nested.length; j++) {
        result.push(nested[j]);
      }
    } else {
      result.push(arr[i]);
    }
  }
  return result;
}

// VIOLATION: @typescript-eslint/no-explicit-any — uses `any` type
export function debounce(fn: any, delay: number): any {
  // VIOLATION: prefer-const
  let timeoutId: any;
  return function (...args: any[]) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => fn(...args), delay);
  };
}

// VIOLATION: @typescript-eslint/no-unused-vars — parameter `options` is not used
export function createSlug(text: string, options: any = {}): string {
  // VIOLATION: prefer-const
  let slug = text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, "")
    .replace(/[\s_-]+/g, "-")
    .replace(/^-+|-+$/g, "");
  return slug;
}

// VIOLATION: @typescript-eslint/no-explicit-any
export function groupBy(items: any[], key: string): Record<string, any[]> {
  // VIOLATION: prefer-const
  let groups: Record<string, any[]> = {};
  for (let i = 0; i < items.length; i++) {
    // VIOLATION: prefer-const
    let groupKey = String(items[i][key]);
    if (!groups[groupKey]) {
      groups[groupKey] = [];
    }
    groups[groupKey].push(items[i]);
  }
  return groups;
}

// VIOLATION: @typescript-eslint/no-unused-vars — `retries` is declared but never read
export function retry(fn: () => Promise<any>, maxAttempts: number): Promise<any> {
  let retries = 0;
  return new Promise(async (resolve, reject) => {
    for (let attempt = 1; attempt <= maxAttempts; attempt++) {
      try {
        // VIOLATION: prefer-const
        let result = await fn();
        resolve(result);
        return;
      } catch (err) {
        if (attempt === maxAttempts) {
          reject(err);
        }
      }
    }
  });
}

// VIOLATION: prefer-const, no-unused-vars
export function getEnvironmentConfig() {
  let env = process.env.NODE_ENV || "development";
  let port = process.env.PORT || "3000";
  // VIOLATION: @typescript-eslint/no-unused-vars
  let secret = process.env.SECRET_KEY || "default-secret";
  // VIOLATION: @typescript-eslint/no-unused-vars
  let dbHost = process.env.DB_HOST || "localhost";

  return {
    env: env,
    port: parseInt(port, 10),
    isProduction: env === "production",
    isDevelopment: env === "development",
  };
}

// VIOLATION: @typescript-eslint/no-explicit-any
export function safeJsonParse(text: string): any {
  try {
    return JSON.parse(text);
  } catch {
    return null;
  }
}

// VIOLATION: @typescript-eslint/no-unused-vars — entire function is unused
function generateId(): string {
  let chars = "abcdefghijklmnopqrstuvwxyz0123456789";
  let result = "";
  for (let i = 0; i < 12; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

// VIOLATION: @typescript-eslint/no-explicit-any
export function deepClone(obj: any): any {
  if (obj === null || typeof obj !== "object") {
    return obj;
  }
  if (Array.isArray(obj)) {
    return obj.map((item: any) => deepClone(item));
  }
  // VIOLATION: prefer-const
  let clone: any = {};
  for (let key in obj) {
    if (obj.hasOwnProperty(key)) {
      clone[key] = deepClone(obj[key]);
    }
  }
  return clone;
}
