import React, { useState } from 'react';
import { DownloadStatus } from '../types';

interface DownloadFormProps {
    onDownloadStart: (download: DownloadStatus) => void;
}

const DownloadForm: React.FC<DownloadFormProps> = ({ onDownloadStart }) => {
    const [url, setUrl] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const response = await fetch('/api/download/start', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url }),
            });
            const data = await response.json();
            onDownloadStart(data);
            setUrl('');
        } catch (error) {
            console.error('Error starting download:', error);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="mb-4">
            <input
                type="text"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                placeholder="Enter video URL"
                className="border p-2 mr-2"
                required
            />
            <button type="submit" className="bg-blue-500 text-white p-2 rounded">
                Start Download
            </button>
        </form>
    );
};

export default DownloadForm;

