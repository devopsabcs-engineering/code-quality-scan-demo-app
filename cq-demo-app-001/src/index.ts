// INTENTIONAL QUALITY VIOLATIONS — DO NOT FIX
// This file contains intentional code quality issues for scanner demonstration:
// - High cyclomatic complexity (CCN > 10)
// - Long functions (> 50 lines)
// - Unused variables
// - Missing type annotations (using 'any')
// - Duplicated code blocks

import express, { Request, Response } from 'express';

const app = express();
const PORT = process.env.PORT || 3000;

app.use(express.json());

// ── VIOLATION: Unused variables ──
let unusedConfig = { debug: true, verbose: false };
let tempData = "this is never used";
var legacyVar = 42;

// ── VIOLATION: High complexity function (CCN > 15) ──
function processOrder(order: any): any {
    let result: any = {};
    let discount = 0;
    let tax = 0;
    let shipping = 0;

    if (order.type === 'premium') {
        if (order.amount > 1000) {
            if (order.loyalty === 'gold') {
                discount = 0.2;
                if (order.coupon) {
                    if (order.coupon === 'SAVE10') {
                        discount += 0.1;
                    } else if (order.coupon === 'SAVE20') {
                        discount += 0.2;
                    } else if (order.coupon === 'SAVE30') {
                        discount += 0.3;
                    }
                }
            } else if (order.loyalty === 'silver') {
                discount = 0.15;
                if (order.coupon) {
                    if (order.coupon === 'SAVE10') {
                        discount += 0.05;
                    } else if (order.coupon === 'SAVE20') {
                        discount += 0.1;
                    }
                }
            } else {
                discount = 0.1;
            }
        } else if (order.amount > 500) {
            discount = 0.05;
            if (order.loyalty === 'gold') {
                discount += 0.1;
            }
        } else {
            discount = 0;
        }
    } else if (order.type === 'standard') {
        if (order.amount > 500) {
            discount = 0.05;
        }
    } else if (order.type === 'basic') {
        discount = 0;
    }

    // Tax calculation — deeply nested
    if (order.region === 'US') {
        if (order.state === 'CA') {
            tax = 0.0725;
        } else if (order.state === 'NY') {
            tax = 0.08;
        } else if (order.state === 'TX') {
            tax = 0.0625;
        } else {
            tax = 0.05;
        }
    } else if (order.region === 'EU') {
        if (order.country === 'DE') {
            tax = 0.19;
        } else if (order.country === 'FR') {
            tax = 0.20;
        } else {
            tax = 0.21;
        }
    } else {
        tax = 0.10;
    }

    // Shipping calculation
    if (order.weight > 50) {
        shipping = 25;
    } else if (order.weight > 20) {
        shipping = 15;
    } else if (order.weight > 5) {
        shipping = 10;
    } else {
        shipping = 5;
    }

    if (order.express) {
        shipping *= 2;
    }

    result.subtotal = order.amount;
    result.discount = order.amount * discount;
    result.tax = (order.amount - result.discount) * tax;
    result.shipping = shipping;
    result.total = order.amount - result.discount + result.tax + result.shipping;

    return result;
}

// ── VIOLATION: Duplicated code block (same pattern in two functions) ──
function formatUserReport(users: any[]): string {
    let report = "=== User Report ===\n";
    for (let i = 0; i < users.length; i++) {
        let user = users[i];
        report += `ID: ${user.id}\n`;
        report += `Name: ${user.name}\n`;
        report += `Email: ${user.email}\n`;
        report += `Status: ${user.status}\n`;
        report += `Created: ${user.createdAt}\n`;
        report += "---\n";
    }
    return report;
}

function formatProductReport(products: any[]): string {
    let report = "=== Product Report ===\n";
    for (let i = 0; i < products.length; i++) {
        let product = products[i];
        report += `ID: ${product.id}\n`;
        report += `Name: ${product.name}\n`;
        report += `Price: ${product.price}\n`;
        report += `Status: ${product.status}\n`;
        report += `Created: ${product.createdAt}\n`;
        report += "---\n";
    }
    return report;
}

// ── VIOLATION: Another high-complexity function ──
function validateInput(data: any): any {
    let errors: string[] = [];

    if (!data.name) {
        errors.push("Name is required");
    } else if (data.name.length < 2) {
        errors.push("Name must be at least 2 characters");
    } else if (data.name.length > 100) {
        errors.push("Name must be less than 100 characters");
    }

    if (!data.email) {
        errors.push("Email is required");
    } else if (!data.email.includes('@')) {
        errors.push("Email must contain @");
    } else if (!data.email.includes('.')) {
        errors.push("Email must contain a domain");
    }

    if (!data.age) {
        errors.push("Age is required");
    } else if (typeof data.age !== 'number') {
        errors.push("Age must be a number");
    } else if (data.age < 0) {
        errors.push("Age must be positive");
    } else if (data.age > 150) {
        errors.push("Age must be realistic");
    }

    if (!data.phone) {
        errors.push("Phone is required");
    } else if (data.phone.length < 10) {
        errors.push("Phone must be at least 10 digits");
    } else if (data.phone.length > 15) {
        errors.push("Phone must be at most 15 digits");
    }

    return { valid: errors.length === 0, errors };
}

// ── Routes ──
app.get('/', (req: Request, res: Response) => {
    res.json({
        app: 'cq-demo-app-001',
        language: 'TypeScript',
        framework: 'Express',
        description: 'Code Quality Demo — intentional violations for scanner testing',
        endpoints: ['/api/order', '/api/users', '/api/products', '/api/validate', '/health']
    });
});

app.post('/api/order', (req: Request, res: Response) => {
    let result = processOrder(req.body);
    res.json(result);
});

app.get('/api/users', (req: Request, res: Response) => {
    let users = [
        { id: 1, name: 'Alice', email: 'alice@example.com', status: 'active', createdAt: '2024-01-01' },
        { id: 2, name: 'Bob', email: 'bob@example.com', status: 'inactive', createdAt: '2024-02-01' }
    ];
    let report = formatUserReport(users);
    res.type('text/plain').send(report);
});

app.get('/api/products', (req: Request, res: Response) => {
    let products = [
        { id: 1, name: 'Widget', price: 9.99, status: 'available', createdAt: '2024-01-15' },
        { id: 2, name: 'Gadget', price: 19.99, status: 'out-of-stock', createdAt: '2024-02-15' }
    ];
    let report = formatProductReport(products);
    res.type('text/plain').send(report);
});

app.post('/api/validate', (req: Request, res: Response) => {
    let result = validateInput(req.body);
    res.json(result);
});

app.get('/health', (req: Request, res: Response) => {
    res.json({ status: 'healthy', timestamp: new Date().toISOString() });
});

app.listen(PORT, () => {
    console.log(`cq-demo-app-001 running on port ${PORT}`);
});

export { app, processOrder, formatUserReport, formatProductReport, validateInput };
