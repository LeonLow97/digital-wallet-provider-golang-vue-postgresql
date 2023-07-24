import { describe, expect, it } from "vitest";
import { currencyConversion } from "../../src/utils/currencyUtils";

describe("Currency Conversion", () => {
  it("should convert amount correctly when currencies are different", () => {
    expect(currencyConversion(100, "USD", "SGD")).toBeCloseTo(135, 2);
    expect(currencyConversion(50, "SGD", "USD")).toBeCloseTo(37, 2);
    expect(currencyConversion(200, "EUR", "SGD")).toBeCloseTo(294, 2);
  });

  it("should return the same amount when currencies are the same", () => {
    expect(currencyConversion(100, "USD", "USD")).toBe(100);
    expect(currencyConversion(50, "SGD", "SGD")).toBe(50);
    expect(currencyConversion(200, "EUR", "EUR")).toBe(200);
  });

  it("should return 0 for invalid inputs", () => {
    expect(currencyConversion(-100, "USD", "SGD")).toBe(0);
    expect(currencyConversion("invalid", "SGD", "USD")).toBe(0);
  });
});
