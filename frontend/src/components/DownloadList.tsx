import React from 'react';
import DownloadItem from './DownloadItem';
import { DownloadStatus } from '../types';

interface DownloadListProps {
    downloads: DownloadStatus[];
    onDownloadUpdate: (download: DownloadStatus) => void;
}

const DownloadList: React.FC<DownloadListProps> = ({ downloads, onDownloadUpdate }) => {
    return (
        <div>
            <h2 className="text-2xl font-bold mb-2">Downloads</h2>
            {downloads.map((download) => (
                <DownloadItem
                    key={download.id}
                    download={download}
                    onDownloadUpdate={onDownloadUpdate}
                />
            ))}
        </div>
    );
};

export default DownloadList;

