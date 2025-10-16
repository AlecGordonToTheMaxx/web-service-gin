'use client';

import { useEffect, useState } from 'react';
import { albumService } from '@/services/albumService';
import type { Album } from '@/types/album';

export default function AlbumDetail({ params }: { params: { id: string } }) {
  const [album, setAlbum] = useState<Album | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchAlbum = async () => {
      try {
        const data = await albumService.getById(Number(params.id));
        setAlbum(data);
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchAlbum();
  }, [params.id]);

  if (loading) return <div className="p-12">Loading...</div>;
  if (!album) return <div className="p-12">Album not found</div>;

  return (
    <div className="min-h-screen p-12">
      <h1 className="text-4xl font-bold mb-4">{album.title}</h1>
      <p className="text-2xl text-purple-600 mb-2">{album.artist}</p>
      <p className="text-3xl font-bold">${album.price.toFixed(2)}</p>
    </div>
  );
}
