using CqDemoApp003.Models;

namespace CqDemoApp003.Services;

/// <summary>
/// Service for managing payroll calculations.
/// INTENTIONAL VIOLATIONS:
/// - ValidatePayrollData() is duplicated from EmployeeService.ValidateEmployeeData()
/// - CalculateSeniorityBonus() is duplicated from EmployeeService
/// - FormatCurrency() is duplicated from EmployeeService
/// - CalculatePayroll() has moderate cyclomatic complexity
/// </summary>
public class PayrollService
{
    private readonly EmployeeService _employeeService;
    private readonly List<PayrollRecord> _payrollRecords = new();
    private int _nextId = 1;

    public PayrollService(EmployeeService employeeService)
    {
        _employeeService = employeeService;
    }

    public List<PayrollRecord> GetAll() => _payrollRecords;

    public PayrollRecord? GetByEmployeeId(int employeeId) =>
        _payrollRecords.FirstOrDefault(p => p.EmployeeId == employeeId);

    // VIOLATION: High Cyclomatic Complexity (CCN > 10)
    // Duplicated validation and calculation patterns from EmployeeService
    public PayrollRecord? CalculatePayroll(int employeeId, string payPeriod, decimal? bonusOverride, decimal? deductionOverride)
    {
        var employee = _employeeService.GetById(employeeId);
        if (employee == null)
            return null;

        if (employee.Status != "Active")
            return null;

        if (string.IsNullOrWhiteSpace(payPeriod))
            return null;

        decimal baseSalary = employee.Salary / 12; // Monthly
        decimal bonus = 0;
        decimal deductions = 0;
        decimal taxRate = 0;

        // VIOLATION: Magic numbers — tax brackets
        if (employee.Salary <= 50000)
            taxRate = 0.15m;
        else if (employee.Salary <= 100000)
            taxRate = 0.22m;
        else if (employee.Salary <= 150000)
            taxRate = 0.28m;
        else if (employee.Salary <= 200000)
            taxRate = 0.33m;
        else
            taxRate = 0.37m;

        // Calculate bonus — duplicated logic from EmployeeService
        if (bonusOverride.HasValue)
        {
            bonus = bonusOverride.Value;
        }
        else
        {
            bonus = CalculateSeniorityBonus(employee.SeniorityLevel, baseSalary, employee.Department);
        }

        // Calculate deductions — department-specific rules
        if (deductionOverride.HasValue)
        {
            deductions = deductionOverride.Value;
        }
        else
        {
            // VIOLATION: Magic numbers — deduction rates
            deductions = baseSalary * 0.06m; // Base deduction
            if (employee.Department == "Engineering")
                deductions += 150.00m; // Tech equipment fund
            else if (employee.Department == "Marketing")
                deductions += 75.00m; // Marketing materials
            else if (employee.Department == "Finance")
                deductions += 100.00m; // Certification fees
        }

        decimal grossPay = baseSalary + bonus;
        decimal taxAmount = grossPay * taxRate;
        decimal netPay = grossPay - taxAmount - deductions;

        var record = new PayrollRecord
        {
            Id = _nextId++,
            EmployeeId = employee.Id,
            EmployeeName = employee.FirstName + " " + employee.LastName,
            Department = employee.Department,
            BaseSalary = baseSalary,
            Bonus = bonus,
            Deductions = deductions,
            TaxAmount = Math.Round(taxAmount, 2),
            NetPay = Math.Round(netPay, 2),
            PayPeriod = payPeriod,
            CalculatedAt = DateTime.Now,
            Status = "Calculated"
        };

        _payrollRecords.Add(record);
        return record;
    }

    // VIOLATION: DUPLICATED from EmployeeService.ValidateEmployeeData()
    // Same validation logic, copied with minor parameter name changes
    public bool ValidatePayrollData(string firstName, string lastName, string email, string department, decimal salary)
    {
        if (string.IsNullOrWhiteSpace(firstName))
            return false;
        if (string.IsNullOrWhiteSpace(lastName))
            return false;
        if (string.IsNullOrWhiteSpace(email) || !email.Contains("@"))
            return false;
        if (string.IsNullOrWhiteSpace(department))
            return false;
        // VIOLATION: Magic numbers
        if (salary < 30000 || salary > 500000)
            return false;
        if (firstName.Length > 100 || lastName.Length > 100)
            return false;
        if (email.Length > 255)
            return false;
        return true;
    }

    // VIOLATION: DUPLICATED from EmployeeService.CalculateSeniorityBonus()
    // Identical logic, copied wholesale
    public decimal CalculateSeniorityBonus(int seniorityLevel, decimal baseSalary, string department)
    {
        decimal bonusPercentage = 0;

        if (seniorityLevel == 1)
            bonusPercentage = 0.02m;
        else if (seniorityLevel == 2)
            bonusPercentage = 0.05m;
        else if (seniorityLevel == 3)
            bonusPercentage = 0.08m;
        else if (seniorityLevel == 4)
            bonusPercentage = 0.12m;
        else if (seniorityLevel >= 5)
            bonusPercentage = 0.15m;

        decimal bonus = baseSalary * bonusPercentage;

        // Department-specific multipliers (magic numbers)
        if (department == "Engineering")
            bonus *= 1.2m;
        else if (department == "Finance")
            bonus *= 1.1m;
        else if (department == "Marketing")
            bonus *= 0.9m;

        return Math.Round(bonus, 2);
    }

    // VIOLATION: DUPLICATED from EmployeeService.FormatCurrency()
    // Identical logic, copied wholesale
    public string FormatCurrency(decimal amount)
    {
        if (amount < 0)
        {
            return "-$" + Math.Abs(amount).ToString("N2");
        }
        else
        {
            return "$" + amount.ToString("N2");
        }
    }

    // VIOLATION: Complex summary method with multiple branches
    public object GeneratePayrollSummary(string payPeriod)
    {
        var records = _payrollRecords.Where(r => r.PayPeriod == payPeriod).ToList();

        if (records.Count == 0)
        {
            return new { PayPeriod = payPeriod, Message = "No payroll records found" };
        }

        decimal totalGross = 0;
        decimal totalNet = 0;
        decimal totalTax = 0;
        decimal totalBonus = 0;

        foreach (var record in records)
        {
            totalGross += record.BaseSalary + record.Bonus;
            totalNet += record.NetPay;
            totalTax += record.TaxAmount;
            totalBonus += record.Bonus;
        }

        return new
        {
            PayPeriod = payPeriod,
            EmployeeCount = records.Count,
            TotalGrossPay = FormatCurrency(totalGross),
            TotalNetPay = FormatCurrency(totalNet),
            TotalTaxWithheld = FormatCurrency(totalTax),
            TotalBonuses = FormatCurrency(totalBonus),
            AverageNetPay = FormatCurrency(totalNet / records.Count)
        };
    }
}
