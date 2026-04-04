// =============================================================================
// calculations.test.ts — INTENTIONAL LOW TEST COVERAGE
// Only tests 1 function (calculateOrderTotal) with a single scenario.
// Leaves processPayment, checkInventoryAvailability, calculateRevenueMetrics,
// and all validator/formatter/helper functions completely untested.
// DO NOT add more tests — this low coverage is intentional for scanner detection.
// =============================================================================

import {
  calculateOrderTotal,
  OrderItem,
  CustomerProfile,
  ShippingAddress,
} from "../src/lib/calculations";

describe("calculateOrderTotal", () => {
  const sampleItems: OrderItem[] = [
    {
      productId: 1,
      name: "Wireless Headphones",
      price: 79.99,
      quantity: 1,
      weight: 0.5,
      category: "Electronics",
    },
    {
      productId: 2,
      name: "Running Shoes",
      price: 129.95,
      quantity: 1,
      weight: 1.2,
      category: "Sports",
    },
  ];

  const sampleCustomer: CustomerProfile = {
    id: "CUST-001",
    tier: "gold",
    region: "west",
    loyaltyPoints: 500,
    memberSince: "2022-01-15",
    isEmployee: false,
  };

  const sampleAddress: ShippingAddress = {
    country: "US",
    state: "CA",
    zipCode: "90210",
    isRemote: false,
  };

  test("should calculate total for a basic order with gold tier discount", () => {
    const result = calculateOrderTotal(
      sampleItems,
      sampleCustomer,
      sampleAddress,
      null,
      false,
      false
    );

    expect(result.subtotal).toBe(209.94);
    expect(result.discount).toBeGreaterThan(0);
    expect(result.tax).toBeGreaterThan(0);
    expect(result.total).toBeGreaterThan(0);
    expect(result.loyaltyPointsEarned).toBeGreaterThan(0);
  });

  test("should return zero discount for bronze tier under $200", () => {
    const bronzeCustomer: CustomerProfile = {
      ...sampleCustomer,
      tier: "bronze",
    };

    const cheapItems: OrderItem[] = [
      {
        productId: 8,
        name: "Water Bottle",
        price: 19.95,
        quantity: 1,
        weight: 0.3,
        category: "Sports",
      },
    ];

    const result = calculateOrderTotal(
      cheapItems,
      bronzeCustomer,
      sampleAddress,
      null,
      false,
      false
    );

    expect(result.subtotal).toBe(19.95);
    expect(result.discount).toBe(0);
  });
});

// NOTE: processPayment is NOT tested — intentional coverage gap
// NOTE: checkInventoryAvailability is NOT tested — intentional coverage gap
// NOTE: calculateRevenueMetrics is NOT tested — intentional coverage gap
// NOTE: All validators.ts functions are NOT tested — intentional coverage gap
// NOTE: All formatters.ts functions are NOT tested — intentional coverage gap
// NOTE: All helpers.ts functions are NOT tested — intentional coverage gap
