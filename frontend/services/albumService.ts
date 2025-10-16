import type { Album, CreateAlbumInput, UpdateAlbumInput } from '@/types/album';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class AlbumServiceError extends Error {
  constructor(
    public status: number,
    message: string
  ) {
    super(message);
    this.name = 'AlbumServiceError';
  }
}

export const albumService = {
  async getAll(): Promise<Album[]> {
    const response = await fetch(`${API_URL}/albums`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      cache: 'no-store', // Disable caching for fresh data
    });

    if (!response.ok) {
      throw new AlbumServiceError(
        response.status,
        `Failed to fetch albums: ${response.statusText}`
      );
    }

    return response.json();
  },

  async getById(id: number): Promise<Album> {
    const response = await fetch(`${API_URL}/albums/${id}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
      cache: 'no-store',
    });

    if (!response.ok) {
      throw new AlbumServiceError(response.status, `Failed to fetch album: ${response.statusText}`);
    }

    return response.json();
  },

  async create(album: CreateAlbumInput): Promise<Album> {
    const response = await fetch(`${API_URL}/albums`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(album),
    });

    if (!response.ok) {
      throw new AlbumServiceError(
        response.status,
        `Failed to create album: ${response.statusText}`
      );
    }

    return response.json();
  },

  async update(id: number, album: UpdateAlbumInput): Promise<Album> {
    const response = await fetch(`${API_URL}/albums/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(album),
    });

    if (!response.ok) {
      throw new AlbumServiceError(
        response.status,
        `Failed to update album: ${response.statusText}`
      );
    }

    return response.json();
  },

  async delete(id: number): Promise<void> {
    const response = await fetch(`${API_URL}/albums/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new AlbumServiceError(
        response.status,
        `Failed to delete album: ${response.statusText}`
      );
    }
  },
};
