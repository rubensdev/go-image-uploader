export interface GlobalConfig {
  lang: string;
  allowedMimetypes: string;
  maxFileSize: number;
  uploadEndpoint: string;
}

declare global {
  interface Window {
    globalConfig: GlobalConfig;
  }
}
