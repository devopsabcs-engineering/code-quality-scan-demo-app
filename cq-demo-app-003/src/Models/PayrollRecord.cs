namespace CqDemoApp003.Models;

/// <summary>
/// Represents a payroll record for an employee.
/// </summary>
public class PayrollRecord
{
    public int Id { get; set; }
    public int EmployeeId { get; set; }
    public string EmployeeName { get; set; } = string.Empty;
    public string Department { get; set; } = string.Empty;
    public decimal BaseSalary { get; set; }
    public decimal Bonus { get; set; }
    public decimal Deductions { get; set; }
    public decimal TaxAmount { get; set; }
    public decimal NetPay { get; set; }
    public string PayPeriod { get; set; } = string.Empty;
    public DateTime CalculatedAt { get; set; }
    public string Status { get; set; } = "Pending";
}

/// <summary>
/// Request model for payroll calculation.
/// </summary>
public class PayrollCalculationRequest
{
    public int EmployeeId { get; set; }
    public string PayPeriod { get; set; } = string.Empty;
    public decimal? BonusOverride { get; set; }
    public decimal? DeductionOverride { get; set; }
}
