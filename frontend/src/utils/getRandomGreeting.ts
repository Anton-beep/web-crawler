const greetings = [
    "Welcome, ",
    "Hello, ",
    "Nice to see you here, ",
    "Greetings, ",
    "Hey, ",
    "Hi, ",
    "Good to see you, ",
    "Hey there, ",
    "Welcome back, ",
    "How are you doing, ",
    "Привет, ",
    "Moin, ",
]

export function getRandomGreeting() {
    return greetings[Math.floor(Math.random() * greetings.length)];
}