import React, { useEffect } from 'react';
import { DownloadStatus } from '../types';

interface DownloadItemProps {
    download: DownloadStatus;
    onDownloadUpdate: (download: DownloadStatus) => void;
}

const DownloadItem: React.FC<DownloadItemProps> = ({ download, onDownloadUpdate }) => {
    useEffect(() => {
        const intervalId = setInterval(async () => {
            if (download.status === 'downloading') {
                try {
                    const response = await fetch(`/api/download/status?id=${download.id}`);
                    const updatedDownload = await response.json();
                    onDownloadUpdate(updatedDownload);
                } catch (error) {
                    console.error('Error fetching download status:', error);
                }
            }
        }, 1000);

        return () => clearInterval(intervalId);
    }, [download.id, download.status, onDownloadUpdate]);

    const handlePauseResume = async () => {
        const action = download.status === 'downloading' ? 'pause' : 'resume';
        try {
            const response = await fetch(`/api/download/${action}?id=${download.id}`);
            const updatedDownload = await response.json();
            onDownloadUpdate(updatedDownload);
        } catch (error) {
            console.error(`Error ${action}ing download:`, error);
        }
    };

    const handleSave = () => {
        window.location.href = `/api/download/save?id=${download.id}`;
    };

    return (
        <div className="border p-4 mb-2">
            <p>URL: {download.url}</p>
            <p>Status: {download.status}</p>
            <p>Progress: {download.progress.toFixed(2)}%</p>
            {download.error && <p className="text-red-500">Error: {download.error}</p>}
            {download.status === 'completed' ? (
                <button onClick={handleSave} className="bg-green-500 text-white p-2 rounded mr-2">
                    Save to Device
                </button>
            ) : (
                <button onClick={handlePauseResume} className="bg-blue-500 text-white p-2 rounded mr-2">
                    {download.status === 'downloading' ? 'Pause' : 'Resume'}
                </button>
            )}
        </div>
    );
};

export default DownloadItem;

