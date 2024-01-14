import formatDateTime from "@/utils/formatDateTime";

describe("formatDateTime", () => {
    it("formats the date time correct in en-US format", () => {
        const dateTime = "2023-08-16T03:33:33.730676Z"

        expect(formatDateTime(dateTime)).toBe("August 16, 2023 at 11:33:33 AM")
    })
})