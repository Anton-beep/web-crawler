import getDomainFromString from "../utils/getDomainFromUrl";

describe('getDomainFromUrl', () => {
    it('should return the domain from a valid URL', () => {
        expect(getDomainFromString('https://www.example.com/path')).toBe('www.example.com');
        expect(getDomainFromString('http://example.com')).toBe('example.com');
        expect(getDomainFromString('https://subdomain.example.co.uk')).toBe('subdomain.example.co.uk');
    });

    it('should return the input string if it is not a valid URL', () => {
        expect(getDomainFromString('not a url')).toBe('not a url');
        expect(getDomainFromString('example')).toBe('example');
        expect(getDomainFromString('')).toBe('');
    });

    it('should handle URLs with ports correctly', () => {
        expect(getDomainFromString('http://localhost:3000')).toBe('localhost');
        expect(getDomainFromString('https://example.com:8080/path')).toBe('example.com');
    });

    it('should handle URLs with different protocols', () => {
        expect(getDomainFromString('ftp://example.com')).toBe('example.com');
    });
});