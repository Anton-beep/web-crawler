export function isUsernameValid(username: string): boolean {
    const minLength = 3;

    return username.length > minLength;
}