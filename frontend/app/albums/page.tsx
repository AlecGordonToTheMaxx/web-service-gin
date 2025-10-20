'use client';

import { useCallback, useEffect, useState } from 'react';
import { albumService } from '@/services/albumService';
import type { Album, CreateAlbumInput } from '@/types/album';

export default function Home() {
  const [albums, setAlbums] = useState<Album[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [form, setForm] = useState<CreateAlbumInput>({
    title: '',
    artist: '',
    price: 0,
  });

  const loadAlbums = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await albumService.getAll();
      setAlbums(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load albums');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadAlbums();
  }, [loadAlbums]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    try {
      if (editingId !== null) {
        await albumService.update(editingId, form);
      } else {
        await albumService.create(form);
      }
      await loadAlbums();
      resetForm();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save album');
    }
  };

  const handleEdit = (album: Album) => {
    setEditingId(album.id);
    setForm({
      title: album.title,
      artist: album.artist,
      price: album.price,
    });
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this album?')) return;

    try {
      setError(null);
      await albumService.delete(id);
      await loadAlbums();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete album');
    }
  };

  const resetForm = () => {
    setForm({ title: '', artist: '', price: 0 });
    setEditingId(null);
  };

  if (loading) {
    return (
      <div
        className="min-h-screen flex items-center justify-center"
        style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
      >
        <div className="text-white text-2xl">Loading albums...</div>
      </div>
    );
  }

  return (
    <div
      className="min-h-screen py-12 px-4 sm:px-6 lg:px-8"
      style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
    >
      <div className="max-w-7xl mx-auto">
        <h1 className="text-5xl font-bold text-white text-center mb-2">Album Manager</h1>
        <p className="text-purple-200 text-center mb-12">Albums ({albums.length})</p>

        {error && (
          <div className="max-w-2xl mx-auto mb-8 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            {error}
          </div>
        )}

        <div className="max-w-2xl mx-auto mb-12 bg-white rounded-2xl shadow-2xl p-8">
          <h2 className="text-2xl font-bold text-gray-800 mb-6">
            {editingId !== null ? 'Edit Album' : 'Add New Album'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
                Title
              </label>
              <input
                id="title"
                name="title"
                type="text"
                required
                value={form.title}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                placeholder="Album title"
              />
            </div>

            <div>
              <label htmlFor="artist" className="block text-sm font-medium text-gray-700 mb-2">
                Artist
              </label>
              <input
                id="artist"
                name="artist"
                type="text"
                required
                value={form.artist}
                onChange={(e) => setForm({ ...form, artist: e.target.value })}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                placeholder="Artist name"
              />
            </div>

            <div>
              <label htmlFor="price" className="block text-sm font-medium text-gray-700 mb-2">
                Price ($)
              </label>
              <input
                id="price"
                name="price"
                type="number"
                step="0.01"
                min="0"
                required
                value={form.price}
                onChange={(e) => setForm({ ...form, price: parseFloat(e.target.value) })}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                placeholder="0.00"
              />
            </div>

            <div className="flex gap-4">
              <button
                type="submit"
                className="flex-1 bg-purple-600 text-white py-3 px-6 rounded-lg hover:bg-purple-700 transition-colors font-semibold"
              >
                {editingId !== null ? 'Update Album' : 'Add Album'}
              </button>
              {editingId !== null && (
                <button
                  type="button"
                  onClick={resetForm}
                  className="px-6 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                >
                  Cancel
                </button>
              )}
            </div>
          </form>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {albums.map((album) => (
            <div
              key={album.id}
              className="bg-gradient-to-br from-gray-50 to-purple-100 rounded-xl shadow-lg p-6 hover:shadow-2xl transition-shadow"
            >
              <h3 className="text-xl font-bold text-gray-800 mb-2">{album.title}</h3>
              <p className="text-purple-700 font-semibold mb-1">{album.artist}</p>
              <p className="text-2xl font-bold text-purple-600 mb-4">${album.price.toFixed(2)}</p>
              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={() => handleEdit(album)}
                  className="flex-1 bg-purple-600 text-white py-2 px-4 rounded-lg hover:bg-purple-700 transition-colors"
                >
                  Edit
                </button>
                <button
                  type="button"
                  onClick={() => handleDelete(album.id)}
                  className="flex-1 bg-red-500 text-white py-2 px-4 rounded-lg hover:bg-red-600 transition-colors"
                >
                  Delete
                </button>
              </div>
            </div>
          ))}
        </div>

        {albums.length === 0 && !loading && (
          <div className="text-center text-white text-xl">
            No albums yet. Add your first album above!
          </div>
        )}
      </div>
    </div>
  );
}
