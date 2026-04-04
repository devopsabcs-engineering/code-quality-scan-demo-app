// =============================================================================
// calculations.ts — INTENTIONAL HIGH CYCLOMATIC COMPLEXITY
// Contains functions with CCN > 15 for scanner detection.
// DO NOT refactor — these violations exist for workshop demonstration.
// =============================================================================

export interface OrderItem {
  productId: number;
  name: string;
  price: number;
  quantity: number;
  weight: number;
  category: string;
}

export interface CustomerProfile {
  id: string;
  tier: "bronze" | "silver" | "gold" | "platinum" | "diamond";
  region: string;
  loyaltyPoints: number;
  memberSince: string;
  isEmployee: boolean;
}

export interface ShippingAddress {
  country: string;
  state: string;
  zipCode: string;
  isRemote: boolean;
}

export interface OrderResult {
  subtotal: number;
  discount: number;
  tax: number;
  shipping: number;
  total: number;
  loyaltyPointsEarned: number;
  estimatedDeliveryDays: number;
}

// HIGH COMPLEXITY: CCN > 20 — deeply nested discount/tax/shipping logic
// This function intentionally has excessive branching for scanner detection.
export function calculateOrderTotal(
  items: OrderItem[],
  customer: CustomerProfile,
  shippingAddress: ShippingAddress,
  couponCode: string | null,
  isExpressShipping: boolean,
  giftWrap: boolean
): OrderResult {
  let subtotal = 0;
  let discount = 0;
  let tax = 0;
  let shipping = 0;
  let loyaltyPointsEarned = 0;
  let estimatedDeliveryDays = 5;

  // Calculate subtotal
  for (let i = 0; i < items.length; i++) {
    subtotal += items[i].price * items[i].quantity;
  }

  // Tier-based discount calculation — deeply nested
  if (customer.tier === "diamond") {
    if (subtotal > 500) {
      discount = subtotal * 0.20;
    } else if (subtotal > 200) {
      discount = subtotal * 0.15;
    } else if (subtotal > 100) {
      discount = subtotal * 0.10;
    } else {
      discount = subtotal * 0.05;
    }
  } else if (customer.tier === "platinum") {
    if (subtotal > 500) {
      discount = subtotal * 0.15;
    } else if (subtotal > 200) {
      discount = subtotal * 0.12;
    } else if (subtotal > 100) {
      discount = subtotal * 0.08;
    } else {
      discount = subtotal * 0.04;
    }
  } else if (customer.tier === "gold") {
    if (subtotal > 500) {
      discount = subtotal * 0.12;
    } else if (subtotal > 200) {
      discount = subtotal * 0.08;
    } else if (subtotal > 100) {
      discount = subtotal * 0.05;
    } else {
      discount = subtotal * 0.02;
    }
  } else if (customer.tier === "silver") {
    if (subtotal > 500) {
      discount = subtotal * 0.08;
    } else if (subtotal > 200) {
      discount = subtotal * 0.05;
    } else {
      discount = subtotal * 0.02;
    }
  } else {
    if (subtotal > 500) {
      discount = subtotal * 0.05;
    } else if (subtotal > 200) {
      discount = subtotal * 0.03;
    } else {
      discount = 0;
    }
  }

  // Coupon code processing
  if (couponCode !== null) {
    if (couponCode === "SAVE10") {
      discount += subtotal * 0.10;
    } else if (couponCode === "SAVE20") {
      discount += subtotal * 0.20;
    } else if (couponCode === "FLAT50") {
      discount += 50;
    } else if (couponCode === "FLAT100") {
      if (subtotal > 200) {
        discount += 100;
      } else {
        discount += 25;
      }
    } else if (couponCode === "EMPLOYEE") {
      if (customer.isEmployee) {
        discount += subtotal * 0.30;
      }
    } else if (couponCode === "LOYALTY") {
      if (customer.loyaltyPoints > 1000) {
        discount += subtotal * 0.15;
      } else if (customer.loyaltyPoints > 500) {
        discount += subtotal * 0.10;
      } else {
        discount += subtotal * 0.05;
      }
    }
  }

  // Ensure discount doesn't exceed subtotal
  if (discount > subtotal) {
    discount = subtotal;
  }

  // Tax calculation by region — another deeply nested block
  let taxableAmount = subtotal - discount;
  if (shippingAddress.country === "US") {
    if (shippingAddress.state === "CA") {
      tax = taxableAmount * 0.0725;
    } else if (shippingAddress.state === "TX") {
      tax = taxableAmount * 0.0625;
    } else if (shippingAddress.state === "NY") {
      tax = taxableAmount * 0.08;
    } else if (shippingAddress.state === "FL") {
      tax = taxableAmount * 0.06;
    } else if (shippingAddress.state === "WA") {
      tax = taxableAmount * 0.065;
    } else if (shippingAddress.state === "OR" || shippingAddress.state === "MT" || shippingAddress.state === "NH") {
      tax = 0;
    } else {
      tax = taxableAmount * 0.07;
    }
  } else if (shippingAddress.country === "CA") {
    if (shippingAddress.state === "ON") {
      tax = taxableAmount * 0.13;
    } else if (shippingAddress.state === "QC") {
      tax = taxableAmount * 0.14975;
    } else if (shippingAddress.state === "AB") {
      tax = taxableAmount * 0.05;
    } else if (shippingAddress.state === "BC") {
      tax = taxableAmount * 0.12;
    } else {
      tax = taxableAmount * 0.15;
    }
  } else if (shippingAddress.country === "UK") {
    tax = taxableAmount * 0.20;
  } else if (shippingAddress.country === "DE" || shippingAddress.country === "FR") {
    tax = taxableAmount * 0.19;
  } else if (shippingAddress.country === "JP") {
    tax = taxableAmount * 0.10;
  } else {
    tax = taxableAmount * 0.15;
  }

  // Shipping calculation — weight-based with zone modifiers
  let totalWeight = 0;
  for (let i = 0; i < items.length; i++) {
    totalWeight += items[i].weight * items[i].quantity;
  }

  if (isExpressShipping) {
    if (totalWeight < 1) {
      shipping = 15.99;
      estimatedDeliveryDays = 1;
    } else if (totalWeight < 5) {
      shipping = 24.99;
      estimatedDeliveryDays = 1;
    } else if (totalWeight < 20) {
      shipping = 39.99;
      estimatedDeliveryDays = 2;
    } else {
      shipping = 59.99;
      estimatedDeliveryDays = 2;
    }
  } else {
    if (totalWeight < 1) {
      shipping = 5.99;
      estimatedDeliveryDays = 5;
    } else if (totalWeight < 5) {
      shipping = 9.99;
      estimatedDeliveryDays = 5;
    } else if (totalWeight < 20) {
      shipping = 14.99;
      estimatedDeliveryDays = 7;
    } else {
      shipping = 24.99;
      estimatedDeliveryDays = 10;
    }
  }

  // Remote area surcharge
  if (shippingAddress.isRemote) {
    shipping *= 1.5;
    estimatedDeliveryDays += 3;
  }

  // Free shipping threshold
  if (subtotal - discount > 150 && !isExpressShipping) {
    shipping = 0;
  }

  // Gift wrap
  if (giftWrap) {
    shipping += items.length * 3.99;
  }

  // Loyalty points earned
  if (customer.tier === "diamond") {
    loyaltyPointsEarned = Math.floor((subtotal - discount) * 5);
  } else if (customer.tier === "platinum") {
    loyaltyPointsEarned = Math.floor((subtotal - discount) * 4);
  } else if (customer.tier === "gold") {
    loyaltyPointsEarned = Math.floor((subtotal - discount) * 3);
  } else if (customer.tier === "silver") {
    loyaltyPointsEarned = Math.floor((subtotal - discount) * 2);
  } else {
    loyaltyPointsEarned = Math.floor((subtotal - discount) * 1);
  }

  let total = (subtotal - discount) + tax + shipping;
  total = Math.round(total * 100) / 100;

  return {
    subtotal,
    discount: Math.round(discount * 100) / 100,
    tax: Math.round(tax * 100) / 100,
    shipping: Math.round(shipping * 100) / 100,
    total,
    loyaltyPointsEarned,
    estimatedDeliveryDays,
  };
}

// HIGH COMPLEXITY: CCN > 15 — payment processing with multiple branches
// Another intentionally complex function for scanner detection.
export function processPayment(
  amount: number,
  paymentMethod: string,
  currency: string,
  customerCountry: string,
  cardType: string | null,
  installments: number,
  isRecurring: boolean
): { success: boolean; transactionId: string; fee: number; message: string } {
  let fee = 0;
  let message = "";
  let success = true;
  let transactionId = `TXN-${Date.now()}`;

  // Payment method processing
  if (paymentMethod === "credit_card") {
    if (cardType === "visa") {
      fee = amount * 0.029 + 0.30;
      if (currency !== "USD") {
        fee += amount * 0.01;
      }
    } else if (cardType === "mastercard") {
      fee = amount * 0.029 + 0.30;
      if (currency !== "USD") {
        fee += amount * 0.01;
      }
    } else if (cardType === "amex") {
      fee = amount * 0.035 + 0.30;
      if (currency !== "USD") {
        fee += amount * 0.02;
      }
    } else if (cardType === "discover") {
      fee = amount * 0.029 + 0.30;
    } else {
      fee = amount * 0.04 + 0.50;
      message = "Unknown card type — surcharge applied";
    }

    if (installments > 1) {
      if (installments <= 3) {
        fee += amount * 0.02;
      } else if (installments <= 6) {
        fee += amount * 0.04;
      } else if (installments <= 12) {
        fee += amount * 0.06;
      } else {
        success = false;
        message = "Maximum 12 installments allowed";
      }
    }
  } else if (paymentMethod === "paypal") {
    fee = amount * 0.034 + 0.49;
    if (customerCountry !== "US" && customerCountry !== "CA") {
      fee += amount * 0.015;
    }
    if (isRecurring) {
      fee += amount * 0.01;
    }
  } else if (paymentMethod === "bank_transfer") {
    if (amount < 100) {
      fee = 5.00;
    } else if (amount < 1000) {
      fee = 10.00;
    } else {
      fee = 25.00;
    }
    if (customerCountry !== "US") {
      fee += 15.00;
    }
  } else if (paymentMethod === "crypto") {
    fee = amount * 0.01;
    if (amount < 50) {
      success = false;
      message = "Minimum crypto payment is $50";
    }
  } else if (paymentMethod === "apple_pay" || paymentMethod === "google_pay") {
    fee = amount * 0.025 + 0.10;
  } else {
    success = false;
    message = "Unsupported payment method";
  }

  // Amount validation
  if (amount <= 0) {
    success = false;
    message = "Invalid payment amount";
  } else if (amount > 50000) {
    success = false;
    message = "Amount exceeds maximum single transaction limit";
  }

  // Currency surcharge
  if (currency === "EUR") {
    fee += amount * 0.005;
  } else if (currency === "GBP") {
    fee += amount * 0.005;
  } else if (currency === "JPY") {
    fee += amount * 0.01;
  } else if (currency !== "USD" && currency !== "CAD") {
    fee += amount * 0.02;
  }

  fee = Math.round(fee * 100) / 100;

  if (!message && success) {
    message = "Payment processed successfully";
  }

  return { success, transactionId, fee, message };
}

// MEDIUM COMPLEXITY: inventory check with nested conditions
export function checkInventoryAvailability(
  requestedItems: { productId: number; quantity: number }[],
  warehouse: string,
  allowBackorder: boolean
): { available: boolean; unavailableItems: number[]; estimatedRestock: string | null } {
  let unavailableItems: number[] = [];
  let estimatedRestock: string | null = null;

  // Simulated inventory data
  const inventory: Record<string, Record<number, number>> = {
    "west": { 1: 50, 2: 30, 3: 100, 4: 0, 5: 45, 6: 200, 7: 35, 8: 150 },
    "east": { 1: 70, 2: 20, 3: 80, 4: 10, 5: 30, 6: 90, 7: 25, 8: 200 },
    "central": { 1: 30, 2: 40, 3: 60, 4: 5, 5: 50, 6: 120, 7: 40, 8: 180 },
  };

  const warehouseStock = inventory[warehouse];
  if (!warehouseStock) {
    return { available: false, unavailableItems: requestedItems.map(i => i.productId), estimatedRestock: null };
  }

  for (let i = 0; i < requestedItems.length; i++) {
    const item = requestedItems[i];
    const stock = warehouseStock[item.productId];

    if (stock === undefined) {
      unavailableItems.push(item.productId);
    } else if (stock < item.quantity) {
      if (allowBackorder) {
        estimatedRestock = "2024-12-25";
      } else {
        unavailableItems.push(item.productId);
      }
    }
  }

  return {
    available: unavailableItems.length === 0,
    unavailableItems,
    estimatedRestock,
  };
}

// Calculate revenue metrics — additional function for coverage gap
export function calculateRevenueMetrics(
  orders: { amount: number; status: string; date: string; region: string }[]
): {
  totalRevenue: number;
  avgOrderValue: number;
  cancelRate: number;
  topRegion: string;
  monthlyBreakdown: Record<string, number>;
} {
  let totalRevenue = 0;
  let cancelledCount = 0;
  let regionRevenue: Record<string, number> = {};
  let monthlyBreakdown: Record<string, number> = {};

  for (let i = 0; i < orders.length; i++) {
    const order = orders[i];

    if (order.status !== "cancelled") {
      totalRevenue += order.amount;

      if (!regionRevenue[order.region]) {
        regionRevenue[order.region] = 0;
      }
      regionRevenue[order.region] += order.amount;

      const month = order.date.substring(0, 7);
      if (!monthlyBreakdown[month]) {
        monthlyBreakdown[month] = 0;
      }
      monthlyBreakdown[month] += order.amount;
    } else {
      cancelledCount++;
    }
  }

  let topRegion = "";
  let maxRevenue = 0;
  for (const region in regionRevenue) {
    if (regionRevenue[region] > maxRevenue) {
      maxRevenue = regionRevenue[region];
      topRegion = region;
    }
  }

  return {
    totalRevenue: Math.round(totalRevenue * 100) / 100,
    avgOrderValue: orders.length > 0 ? Math.round((totalRevenue / orders.length) * 100) / 100 : 0,
    cancelRate: orders.length > 0 ? Math.round((cancelledCount / orders.length) * 100) / 100 : 0,
    topRegion: topRegion || "N/A",
    monthlyBreakdown,
  };
}
