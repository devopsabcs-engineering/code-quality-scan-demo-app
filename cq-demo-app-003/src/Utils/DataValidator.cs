using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.RegularExpressions;
using System.Net.Http;
using System.Threading.Tasks;

namespace CqDemoApp003.Utils;

/// <summary>
/// Data validation utility class.
/// INTENTIONAL VIOLATIONS:
/// - Non-standard naming (methods starting with lowercase)
/// - Magic numbers without constants
/// - Methods with too many parameters (> 7)
/// - Empty catch blocks
/// - Inconsistent naming conventions
/// </summary>
public class DataValidator
{
    // VIOLATION: Public field instead of property (CA1051)
    public int validationCount = 0;

    // VIOLATION: Public field instead of property (CA1051)
    public string lastError = "";

    // VIOLATION: Method name starts with lowercase (CA1707, naming convention)
    public bool validateEmail(string email)
    {
        validationCount++;
        // VIOLATION: Empty catch block (CA1031)
        try
        {
            if (string.IsNullOrEmpty(email))
                return false;

            // VIOLATION: Magic number — minimum email length
            if (email.Length < 5)
                return false;

            // VIOLATION: Magic number — maximum email length
            if (email.Length > 320)
                return false;

            return email.Contains("@") && email.Contains(".");
        }
        catch (Exception)
        {
            // VIOLATION: Empty catch block — swallowing exceptions
            return false;
        }
    }

    // VIOLATION: Method name starts with lowercase (naming convention)
    public bool validatePhoneNumber(string phone)
    {
        validationCount++;
        // VIOLATION: Empty catch block
        try
        {
            if (string.IsNullOrEmpty(phone))
                return false;

            var digitsOnly = new string(phone.Where(char.IsDigit).ToArray());

            // VIOLATION: Magic numbers — phone number length validation
            if (digitsOnly.Length < 7 || digitsOnly.Length > 15)
                return false;

            return true;
        }
        catch (Exception)
        {
            // VIOLATION: Empty catch block — swallowing exceptions
            return false;
        }
    }

    // VIOLATION: Method name starts with lowercase (naming convention)
    public bool validateSalary(decimal salary)
    {
        validationCount++;
        // VIOLATION: Magic numbers — salary bounds
        if (salary < 15000)
        {
            lastError = "Salary below minimum";
            return false;
        }
        if (salary > 999999)
        {
            lastError = "Salary above maximum";
            return false;
        }
        return true;
    }

    // VIOLATION: Too many parameters (> 7) — CA1026
    // VIOLATION: Method name starts with lowercase (naming convention)
    public bool validateCompleteRecord(
        string firstName,
        string lastName,
        string email,
        string phone,
        string department,
        string title,
        decimal salary,
        int seniorityLevel,
        string managerId,
        string status)
    {
        validationCount++;

        if (string.IsNullOrWhiteSpace(firstName)) return false;
        if (string.IsNullOrWhiteSpace(lastName)) return false;
        if (!validateEmail(email)) return false;
        if (!validatePhoneNumber(phone)) return false;
        if (string.IsNullOrWhiteSpace(department)) return false;
        if (string.IsNullOrWhiteSpace(title)) return false;
        if (!validateSalary(salary)) return false;

        // VIOLATION: Magic numbers — seniority bounds
        if (seniorityLevel < 1 || seniorityLevel > 10)
            return false;

        if (string.IsNullOrWhiteSpace(status))
            return false;

        // VIOLATION: Magic number — allowed statuses hardcoded
        var validStatuses = new[] { "Active", "Suspended", "Terminated", "OnLeave" };
        if (!validStatuses.Contains(status))
            return false;

        return true;
    }

    // VIOLATION: Method name starts with lowercase (naming convention)
    // VIOLATION: High complexity with nested conditions
    public string validateAndCategorizeDepartment(string department, int headCount, decimal budget)
    {
        validationCount++;

        if (string.IsNullOrWhiteSpace(department))
            return "INVALID";

        // VIOLATION: Magic numbers throughout
        if (department == "Engineering")
        {
            if (headCount > 200)
                return "LARGE_ENGINEERING";
            else if (headCount > 50)
                return "MEDIUM_ENGINEERING";
            else
                return "SMALL_ENGINEERING";
        }
        else if (department == "Finance")
        {
            if (budget > 5000000)
                return "HIGH_BUDGET_FINANCE";
            else if (budget > 1000000)
                return "STANDARD_FINANCE";
            else
                return "LOW_BUDGET_FINANCE";
        }
        else if (department == "Marketing")
        {
            if (headCount > 100 && budget > 2000000)
                return "MAJOR_MARKETING";
            else if (headCount > 30)
                return "STANDARD_MARKETING";
            else
                return "SMALL_MARKETING";
        }
        else if (department == "HR")
        {
            if (headCount > 50)
                return "LARGE_HR";
            else
                return "STANDARD_HR";
        }
        else
        {
            return "OTHER_" + department.ToUpper();
        }
    }

    // VIOLATION: Method name starts with lowercase (naming convention)
    // VIOLATION: Empty catch block
    public decimal parseAmount(string amountString)
    {
        validationCount++;
        try
        {
            // Strip currency symbols and whitespace
            var cleaned = amountString.Replace("$", "").Replace(",", "").Trim();
            return decimal.Parse(cleaned);
        }
        catch (Exception)
        {
            // VIOLATION: Empty catch block — returns 0 silently
            return 0;
        }
    }
}
