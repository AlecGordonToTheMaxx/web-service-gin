import React, { useState, useEffect } from 'react';
import { albumService } from './services/albumService';
import './App.css';

function App() {
  const [albums, setAlbums] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [form, setForm] = useState({ title: '', artist: '', price: '' });
  const [editingId, setEditingId] = useState(null);

  useEffect(() => {
    loadAlbums();
  }, []);

  const loadAlbums = async () => {
    try {
      setLoading(true);
      const data = await albumService.getAll();
      setAlbums(data || []);
      setError(null);
    } catch (err) {
      setError('Failed to load albums');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const albumData = {
        title: form.title,
        artist: form.artist,
        price: parseFloat(form.price),
      };

      if (editingId) {
        await albumService.update(editingId, albumData);
        setEditingId(null);
      } else {
        await albumService.create(albumData);
      }

      setForm({ title: '', artist: '', price: '' });
      await loadAlbums();
    } catch (err) {
      setError(editingId ? 'Failed to update album' : 'Failed to create album');
      console.error(err);
    }
  };

  const handleEdit = (album) => {
    setForm({
      title: album.title,
      artist: album.artist,
      price: album.price.toString(),
    });
    setEditingId(album.id);
  };

  const handleCancelEdit = () => {
    setForm({ title: '', artist: '', price: '' });
    setEditingId(null);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this album?')) {
      return;
    }

    try {
      await albumService.delete(id);
      await loadAlbums();
    } catch (err) {
      setError('Failed to delete album');
      console.error(err);
    }
  };

  if (loading) {
    return <div className="App"><div className="loading">Loading...</div></div>;
  }

  return (
    <div className="App">
      <header className="App-header">
        <h1>Album Manager</h1>
      </header>

      {error && <div className="error-message">{error}</div>}

      <div className="container">
        <form onSubmit={handleSubmit} className="album-form">
          <h2>{editingId ? 'Edit Album' : 'Add New Album'}</h2>
          <input
            type="text"
            placeholder="Title"
            value={form.title}
            onChange={(e) => setForm({ ...form, title: e.target.value })}
            required
          />
          <input
            type="text"
            placeholder="Artist"
            value={form.artist}
            onChange={(e) => setForm({ ...form, artist: e.target.value })}
            required
          />
          <input
            type="number"
            step="0.01"
            placeholder="Price"
            value={form.price}
            onChange={(e) => setForm({ ...form, price: e.target.value })}
            required
          />
          <div className="form-buttons">
            <button type="submit" className="btn-primary">
              {editingId ? 'Update Album' : 'Add Album'}
            </button>
            {editingId && (
              <button type="button" onClick={handleCancelEdit} className="btn-secondary">
                Cancel
              </button>
            )}
          </div>
        </form>

        <div className="albums-list">
          <h2>Albums ({albums.length})</h2>
          {albums.length === 0 ? (
            <p className="no-albums">No albums yet. Add one above!</p>
          ) : (
            <div className="albums-grid">
              {albums.map((album) => (
                <div key={album.id} className="album-card">
                  <h3>{album.title}</h3>
                  <p className="artist">{album.artist}</p>
                  <p className="price">${album.price.toFixed(2)}</p>
                  <div className="card-buttons">
                    <button onClick={() => handleEdit(album)} className="btn-edit">
                      Edit
                    </button>
                    <button onClick={() => handleDelete(album.id)} className="btn-delete">
                      Delete
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default App;
