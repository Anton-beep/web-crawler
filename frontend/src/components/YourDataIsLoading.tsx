import {motion} from "framer-motion";
import {Highlight} from "@/components/ui/hero-highlight.tsx";

export function YourDataIsLoading() {
    return (
        <div
            className="flex flex-row items-center justify-center py-20 h-screen md:h-auto relative w-full">
            <div className="max-w-7xl mx-auto w-full relative overflow-hidden h-full md:h-[4rem] px-4">
                <motion.div
                    initial={{
                        opacity: 0,
                        y: 20,
                    }}
                    animate={{
                        opacity: 1,
                        y: 0,
                    }}
                    transition={{
                        duration: 1,
                    }}
                    className="div"
                >
                    <h2 className="text-center text-white">
                        <Highlight className="text-white">
                            Your data is loading
                        </Highlight>
                    </h2>
                </motion.div>
            </div>
        </div>
    )
}