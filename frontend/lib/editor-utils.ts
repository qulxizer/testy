import runnersData from "@/data/runners.json";
import stubsData from "@/data/stubs.json";
import { Exercise } from "@/lib/types";

const runners: Record<string, string> = runnersData;
const stubs: Record<string, string> = stubsData;

export const getInitialFiles = (exercise: Exercise): Record<string, string> => {
  const filename = exercise.expectedFiles?.[0] || "main.go";

  if (filename.endsWith("main.go")) {
    return {
      [filename]: `package main\n\nimport "github.com/01-edu/z01"\n\nfunc main() {\n\t// Write your code here\n}\n`,
    };
  }

  // Auto-loaded prototype or fallback
  const solutionCode =
    stubs[exercise.key] ||
    `package piscine\n\n// Solution for ${exercise.key}\n`;

  const runnerCode =
    runners[exercise.key] ||
    `package main\n\nimport (\n\t"fmt"\n\t"piscine"\n)\n\nfunc main() {\n\t// Test your ${exercise.name} code here\n}\n`;

  return {
    [filename]: solutionCode,
    "main.go": runnerCode,
  };
};
