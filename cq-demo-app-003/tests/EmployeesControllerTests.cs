using System.Net;
using System.Net.Http.Json;
using Microsoft.AspNetCore.Mvc.Testing;
using Xunit;

namespace CqDemoApp003.Tests;

/// <summary>
/// INTENTIONAL VIOLATION: LOW TEST COVERAGE
/// Only tests the GET /api/employees endpoint.
/// Services (EmployeeService, PayrollService) are completely untested.
/// Utils (DataValidator, ReportGenerator) are completely untested.
/// Payroll endpoints are completely untested.
/// Employee actions (hire, fire, promote, transfer, suspend, reinstate) are untested.
/// Search functionality is untested.
/// </summary>
public class EmployeesControllerTests : IClassFixture<WebApplicationFactory<Program>>
{
    private readonly HttpClient _client;

    public EmployeesControllerTests(WebApplicationFactory<Program> factory)
    {
        _client = factory.CreateClient();
    }

    [Fact]
    public async Task GetAll_ReturnsOk()
    {
        // Arrange & Act
        var response = await _client.GetAsync("/api/employees");

        // Assert
        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
    }

    [Fact]
    public async Task GetAll_ReturnsJsonContent()
    {
        // Arrange & Act
        var response = await _client.GetAsync("/api/employees");
        var content = await response.Content.ReadAsStringAsync();

        // Assert
        Assert.NotNull(content);
        Assert.Contains("success", content);
        Assert.Contains("data", content);
    }

    [Fact]
    public async Task GetById_WithValidId_ReturnsOk()
    {
        // Arrange & Act
        var response = await _client.GetAsync("/api/employees/1");

        // Assert
        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
    }

    [Fact]
    public async Task GetById_WithInvalidId_ReturnsNotFound()
    {
        // Arrange & Act
        var response = await _client.GetAsync("/api/employees/999");

        // Assert
        Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
    }

    [Fact]
    public async Task RootEndpoint_ReturnsOk()
    {
        // Arrange & Act
        var response = await _client.GetAsync("/");

        // Assert
        Assert.Equal(HttpStatusCode.OK, response.StatusCode);
        var content = await response.Content.ReadAsStringAsync();
        Assert.Contains("cq-demo-app-003", content);
    }

    // NOTE: No tests for:
    // - POST /api/employees (Create)
    // - POST /api/employees/{id}/action (ProcessAction)
    // - GET /api/employees/search (SearchEmployees)
    // - GET /api/payroll (GetAll payroll)
    // - POST /api/payroll/calculate (CalculatePayroll)
    // - GET /api/payroll/summary/{payPeriod} (GetSummary)
    // - EmployeeService.ProcessEmployeeAction() — all 6 action types
    // - EmployeeService.ValidateEmployeeData()
    // - EmployeeService.CalculateSeniorityBonus()
    // - PayrollService.CalculatePayroll()
    // - PayrollService.ValidatePayrollData()
    // - DataValidator — all validation methods
    // - ReportGenerator — all report generation methods
}
