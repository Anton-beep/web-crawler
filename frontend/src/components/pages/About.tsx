import AnimationOnScroll from "@/components/AnimationOnScroll.tsx";
import {motion} from "framer-motion";
import {Highlight} from "../ui/hero-highlight";
import {getSampleArcs, globeConfig} from "@/data/globeData.ts";
import {World} from "@/components/ui/globe.tsx";

export default function About() {
    return (
        <div>
            <div
                className="flex flex-row items-center justify-center py-20 h-screen md:h-auto relative w-full">
                <div className="max-w-7xl mx-auto w-full relative overflow-hidden h-full md:h-[40rem] px-4">
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
                        <h2 className="text-center text-xl md:text-4xl font-bold text-white">
                            We crawl around the web to collect data for{" "}
                            <Highlight className="text-white">
                                you
                            </Highlight>
                        </h2>
                    </motion.div>
                    <div
                        className="absolute w-full bottom-0 inset-x-0 h-40 bg-gradient-to-b pointer-events-none select-none from-transparent to-background z-40"/>
                    <div className="absolute w-full -bottom-20 h-72 md:h-full z-10">
                        <World data={getSampleArcs()} globeConfig={globeConfig}/>
                    </div>
                </div>
            </div>
            <div className="flex flex-col md:flex-row items-center gap-8 mb-52 text-primary">
                <div className="md:w-1/2">
                    <img alt={"we need a picture here"} src="about_page_1.png"
                         className="rounded border border-zinc-600"/>
                </div>
                <div className="md:w-1/2">
                    <AnimationOnScroll startState={"transform -translate-x-3/4"}
                                       endState={"transform translate-x-0"}>
                        <p className="text-center text-3xl text-primary font-bold mx-32">We aim to create a robust
                            tool
                            that
                            simplifies and enhances the <span className="text-accent">data collection process, empowering businesses and analysts to make
                    informed decisions.</span> Our crawler operates on a
                            microservice architecture, ensuring <span className="text-accent">high performance through horizontal scalability.</span> Each
                            page is processed by dedicated services for optimal efficiency.</p>
                    </AnimationOnScroll>
                </div>
            </div>

            <div className="flex flex-col md:flex-row items-center gap-8 mb-52 text-primary">
                <div className="md:w-1/2">
                    <AnimationOnScroll startState={"transform translate-x-3/4"}
                                       endState={"transform translate-x-0"}>
                        <p className="text-center text-3xl text-primary font-bold mx-32">Our <span
                            className="text-warning">super intelligent tools </span> <span className="text-zinc-700">(still experimental tho)</span> can help you with solving any complex
                            problems you might encounter on your way.
                        </p>
                    </AnimationOnScroll>
                </div>
                <div className="md:w-1/2">
                    <img alt={"we need a picture here"} src="about_page_2.png"
                         className="rounded border border-zinc-600"/>
                </div>
            </div>
        </div>
    );
}
