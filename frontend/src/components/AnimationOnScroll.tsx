import {ReactNode, useEffect, useRef, useState} from "react";

export default function AnimationOnScroll({children, startState, endState}: { children: ReactNode, startState: string, endState: string }) {
    const [isVisible, setIsVisible] = useState(false);
    const [hasAnimated, setHasAnimated] = useState(false);
    const textRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const observer = new IntersectionObserver(
            ([entry]) => {
                if (entry.isIntersecting && !hasAnimated) {
                    setIsVisible(true);
                    setHasAnimated(true);
                }
            },
            {threshold: 0.1}
        );

        if (textRef.current) {
            observer.observe(textRef.current);
        }

        return () => {
            if (textRef.current) {
                observer.unobserve(textRef.current);
            }
        };
    }, [hasAnimated]);

    // isVisible ? "transform translate-x-0" : "transform -translate-x-2/4"

    return (
        <div
            ref={textRef}
            className={`text-4xl font-bold transition-transform duration-1000 ${
                isVisible ? endState : startState
            }`}
        >
            {children}
        </div>
    );
}