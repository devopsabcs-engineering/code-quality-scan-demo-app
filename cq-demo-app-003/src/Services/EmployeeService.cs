using CqDemoApp003.Models;

namespace CqDemoApp003.Services;

/// <summary>
/// Service for managing employees.
/// INTENTIONAL VIOLATIONS:
/// - ProcessEmployeeAction() has CCN > 15 (high cyclomatic complexity)
/// - ValidateEmployeeData() duplicated in PayrollService
/// - CalculateSeniorityBonus() duplicated in PayrollService
/// - Magic numbers throughout
/// </summary>
public class EmployeeService
{
    private readonly List<Employee> _employees = new()
    {
        new Employee { Id = 1, FirstName = "Alice", LastName = "Johnson", Email = "alice@example.com", Department = "Engineering", Title = "Senior Developer", Salary = 120000, HireDate = DateTime.Now.AddYears(-5), Status = "Active", SeniorityLevel = 3 },
        new Employee { Id = 2, FirstName = "Bob", LastName = "Smith", Email = "bob@example.com", Department = "Marketing", Title = "Marketing Manager", Salary = 95000, HireDate = DateTime.Now.AddYears(-3), Status = "Active", SeniorityLevel = 2 },
        new Employee { Id = 3, FirstName = "Carol", LastName = "Williams", Email = "carol@example.com", Department = "Finance", Title = "Financial Analyst", Salary = 85000, HireDate = DateTime.Now.AddYears(-1), Status = "Active", SeniorityLevel = 1 },
        new Employee { Id = 4, FirstName = "David", LastName = "Brown", Email = "david@example.com", Department = "Engineering", Title = "Junior Developer", Salary = 75000, HireDate = DateTime.Now.AddMonths(-6), Status = "Active", SeniorityLevel = 1 },
        new Employee { Id = 5, FirstName = "Eve", LastName = "Davis", Email = "eve@example.com", Department = "HR", Title = "HR Specialist", Salary = 70000, HireDate = DateTime.Now.AddYears(-2), Status = "Active", SeniorityLevel = 2 }
    };

    private int _nextId = 6;

    public List<Employee> GetAll() => _employees;

    public Employee? GetById(int id) => _employees.FirstOrDefault(e => e.Id == id);

    public Employee Create(Employee employee)
    {
        employee.Id = _nextId++;
        employee.HireDate = DateTime.Now;
        employee.Status = "Active";
        _employees.Add(employee);
        return employee;
    }

    // VIOLATION: High Cyclomatic Complexity (CCN > 15)
    // Nested switch with nested if/else for department rules, seniority, budget
    public string ProcessEmployeeAction(int employeeId, string actionType, string targetDepartment, decimal? newSalary, string reason)
    {
        var employee = _employees.FirstOrDefault(e => e.Id == employeeId);
        if (employee == null)
        {
            return "Employee not found";
        }

        string result = "";

        switch (actionType.ToLower())
        {
            case "hire":
                if (employee.Status == "Active")
                {
                    result = "Employee is already active";
                }
                else if (employee.Status == "Terminated")
                {
                    if (employee.SeniorityLevel > 2)
                    {
                        result = "Rehiring senior employees requires VP approval";
                    }
                    else
                    {
                        employee.Status = "Active";
                        employee.HireDate = DateTime.Now;
                        result = "Employee rehired successfully";
                    }
                }
                else
                {
                    employee.Status = "Active";
                    result = "Employee hired";
                }
                break;

            case "fire":
                if (employee.Status != "Active")
                {
                    result = "Cannot fire inactive employee";
                }
                else if (employee.Department == "Engineering" && employee.SeniorityLevel >= 3)
                {
                    if (string.IsNullOrEmpty(reason))
                    {
                        result = "Reason required for terminating senior engineers";
                    }
                    else if (reason.Length < 50)
                    {
                        result = "Detailed reason required (minimum 50 characters) for senior engineers";
                    }
                    else
                    {
                        employee.Status = "Terminated";
                        result = "Senior engineer terminated with documented reason";
                    }
                }
                else if (employee.Department == "Finance")
                {
                    if (DateTime.Now.Month == 12 || DateTime.Now.Month == 1)
                    {
                        result = "Cannot terminate Finance employees during fiscal year-end";
                    }
                    else
                    {
                        employee.Status = "Terminated";
                        result = "Finance employee terminated";
                    }
                }
                else
                {
                    employee.Status = "Terminated";
                    result = "Employee terminated";
                }
                break;

            case "promote":
                if (employee.Status != "Active")
                {
                    result = "Cannot promote inactive employee";
                }
                else if (employee.SeniorityLevel >= 5)
                {
                    result = "Employee is already at maximum seniority level";
                }
                else if (employee.SeniorityLevel < 2 && (DateTime.Now - employee.HireDate).TotalDays < 365)
                {
                    result = "Junior employees must complete at least 1 year before promotion";
                }
                else
                {
                    employee.SeniorityLevel++;
                    if (newSalary.HasValue && newSalary.Value > employee.Salary)
                    {
                        // VIOLATION: Magic number — max raise percentage
                        if (newSalary.Value > employee.Salary * 1.30m)
                        {
                            result = "Salary increase cannot exceed 30% in a single promotion";
                        }
                        else
                        {
                            employee.Salary = newSalary.Value;
                            result = $"Employee promoted to level {employee.SeniorityLevel} with new salary";
                        }
                    }
                    else
                    {
                        // VIOLATION: Magic number — default raise
                        employee.Salary *= 1.10m;
                        result = $"Employee promoted to level {employee.SeniorityLevel} with 10% raise";
                    }
                }
                break;

            case "transfer":
                if (employee.Status != "Active")
                {
                    result = "Cannot transfer inactive employee";
                }
                else if (string.IsNullOrEmpty(targetDepartment))
                {
                    result = "Target department is required for transfers";
                }
                else if (targetDepartment == employee.Department)
                {
                    result = "Employee is already in the target department";
                }
                else if (targetDepartment == "Engineering" && employee.SeniorityLevel < 2)
                {
                    result = "Engineering department requires seniority level 2 or higher";
                }
                else if (targetDepartment == "Finance" && employee.Department == "Marketing")
                {
                    // VIOLATION: Magic number — waiting period in days
                    if ((DateTime.Now - employee.HireDate).TotalDays < 180)
                    {
                        result = "Marketing to Finance transfers require 180 days tenure";
                    }
                    else
                    {
                        employee.Department = targetDepartment;
                        result = "Employee transferred from Marketing to Finance";
                    }
                }
                else
                {
                    employee.Department = targetDepartment;
                    result = $"Employee transferred to {targetDepartment}";
                }
                break;

            case "suspend":
                if (employee.Status == "Suspended")
                {
                    result = "Employee is already suspended";
                }
                else if (employee.Status != "Active")
                {
                    result = "Only active employees can be suspended";
                }
                else
                {
                    if (string.IsNullOrEmpty(reason))
                    {
                        result = "Reason is required for suspension";
                    }
                    else
                    {
                        employee.Status = "Suspended";
                        result = "Employee suspended";
                    }
                }
                break;

            case "reinstate":
                if (employee.Status != "Suspended")
                {
                    result = "Only suspended employees can be reinstated";
                }
                else
                {
                    if (employee.Department == "Finance" && employee.SeniorityLevel >= 3)
                    {
                        result = "Senior Finance employees require board approval for reinstatement";
                    }
                    else
                    {
                        employee.Status = "Active";
                        result = "Employee reinstated";
                    }
                }
                break;

            default:
                result = $"Unknown action type: {actionType}";
                break;
        }

        return result;
    }

    // VIOLATION: This method is duplicated in PayrollService (code duplication)
    public bool ValidateEmployeeData(string firstName, string lastName, string email, string department, decimal salary)
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

    // VIOLATION: This method is duplicated in PayrollService (code duplication)
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

    // VIOLATION: This method is duplicated in PayrollService (code duplication)
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
}
