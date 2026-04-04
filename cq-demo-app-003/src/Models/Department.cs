namespace CqDemoApp003.Models;

/// <summary>
/// Represents a department in the organization.
/// </summary>
public class Department
{
    public int Id { get; set; }
    public string Name { get; set; } = string.Empty;
    public string Code { get; set; } = string.Empty;
    public decimal Budget { get; set; }
    public int HeadCount { get; set; }
    public int MaxHeadCount { get; set; }
    public string ManagerName { get; set; } = string.Empty;
    public bool IsActive { get; set; } = true;
}
