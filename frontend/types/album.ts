export interface Album {
  id: number;
  title: string;
  artist: string;
  price: number;
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
}

export interface CreateAlbumInput {
  title: string;
  artist: string;
  price: number;
}

export interface UpdateAlbumInput {
  title: string;
  artist: string;
  price: number;
}
