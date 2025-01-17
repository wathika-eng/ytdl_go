export interface DownloadStatus {
    id: string;
    url: string;
    progress: number;
    status: string;
    error?: string;
    createdAt: string;
}

