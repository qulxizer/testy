export interface Exercise {
  id: number;
  key: string;
  name: string;
  path: string;
  difficulty?: number;
  language: string;
  expectedFiles: string[];
  allowedFunctions?: string[];
  testImage?: string;
}
