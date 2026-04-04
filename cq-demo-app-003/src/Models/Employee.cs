namespace CqDemoApp003.Models;

/// <summary>
/// Represents an employee in the system.
/// </summary>
public class Employee
{
    public int Id { get; set; }
    public string FirstName { get; set; } = string.Empty;
    public string LastName { get; set; } = string.Empty;
    public string Email { get; set; } = string.Empty;
    public string Department { get; set; } = string.Empty;
    public string Title { get; set; } = string.Empty;
    public decimal Salary { get; set; }
    public DateTime HireDate { get; set; }
    public string Status { get; set; } = "Active";
    public int SeniorityLevel { get; set; } = 1;
    public string ManagerId { get; set; } = string.Empty;
}

/// <summary>
/// Request model for employee actions.
/// </summary>
public class EmployeeActionRequest
{
    public string ActionType { get; set; } = string.Empty;
    public string TargetDepartment { get; set; } = string.Empty;
    public decimal? NewSalary { get; set; }
    public string Reason { get; set; } = string.Empty;
}
