using CqDemoApp003.Services;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSingleton<EmployeeService>();
builder.Services.AddSingleton<PayrollService>();

var app = builder.Build();

app.MapControllers();
app.MapGet("/", () => Results.Ok(new
{
    application = "cq-demo-app-003",
    description = "Employee & Payroll Management API",
    framework = "ASP.NET Core 8.0",
    language = "C#",
    endpoints = new[]
    {
        "GET  /api/employees",
        "GET  /api/employees/{id}",
        "POST /api/employees",
        "POST /api/employees/{id}/action",
        "GET  /api/payroll",
        "POST /api/payroll/calculate"
    }
}));

app.Run();

// Make Program class accessible for integration tests
public partial class Program { }
