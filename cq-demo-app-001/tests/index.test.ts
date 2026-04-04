// Intentionally minimal tests — low coverage for scanner demonstration
import { processOrder } from '../src/index';

describe('processOrder', () => {
    it('should calculate basic order total', () => {
        const order = {
            type: 'standard',
            amount: 100,
            region: 'US',
            state: 'CA',
            weight: 3,
            express: false
        };
        const result = processOrder(order);
        expect(result.total).toBeGreaterThan(0);
    });
});

// NOTE: formatUserReport, formatProductReport, and validateInput are NOT tested
// This is intentional to produce low coverage findings
