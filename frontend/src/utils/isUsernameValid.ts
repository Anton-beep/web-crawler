export function isUsernameValid(username: string): boolean {
    console.log(username);
    const minLength = 3;

    return username.length >= minLength;
}