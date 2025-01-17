import React, { useState, useEffect } from 'react'
import DownloadForm from './components/DownloadForm';
import DownloadList from './components/DownloadList';
import { DownloadStatus } from './types';


const App: React.FC = () => {
  const [downloads, setDownloads] = useState<DownloadStatus[]>([]);

  useEffect(() => {
    fetchDownloads();
  }, []);

  const fetchDownloads = async () => {
    try {
      const response = await fetch('/api/downloads');
      const data = await response.json();
      setDownloads(data);
    } catch (error) {
      console.error('Error fetching downloads:', error);
    }
  };

  const addDownload = (newDownload: DownloadStatus) => {
    setDownloads((prevDownloads) => [...prevDownloads, newDownload]);
  };

  const updateDownload = (updatedDownload: DownloadStatus) => {
    setDownloads((prevDownloads) =>
      prevDownloads.map((download) =>
        download.id === updatedDownload.id ? updatedDownload : download
      )
    );
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold mb-4">Video Downloader</h1>
      <DownloadForm onDownloadStart={addDownload} />
      <DownloadList downloads={downloads} onDownloadUpdate={updateDownload} />
    </div>
  );
};

export default App;

