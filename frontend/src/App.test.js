import { render, screen, waitFor } from '@testing-library/react';
import App from './App';
import { albumService } from './services/albumService';

// Mock the albumService
jest.mock('./services/albumService');

describe('App Component', () => {
  beforeEach(() => {
    // Reset mocks before each test
    jest.clearAllMocks();
  });

  test('renders Album Manager header', async () => {
    // Mock the getAll method to return empty array
    albumService.getAll.mockResolvedValue([]);

    render(<App />);

    // Wait for loading to complete
    await waitFor(() => {
      expect(screen.queryByText(/loading/i)).not.toBeInTheDocument();
    });

    // Check if header is rendered
    const headerElement = screen.getByText(/album manager/i);
    expect(headerElement).toBeInTheDocument();
  });

  test('renders Add New Album form', async () => {
    albumService.getAll.mockResolvedValue([]);

    render(<App />);

    await waitFor(() => {
      expect(screen.queryByText(/loading/i)).not.toBeInTheDocument();
    });

    // Check if form is rendered
    const formHeader = screen.getByText(/add new album/i);
    expect(formHeader).toBeInTheDocument();

    // Check for input fields
    expect(screen.getByPlaceholderText(/title/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/artist/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/price/i)).toBeInTheDocument();
  });

  test('displays albums when loaded', async () => {
    const mockAlbums = [
      { id: 1, title: 'The Wall', artist: 'Pink Floyd', price: 24.99 },
      { id: 2, title: 'Dark Side of the Moon', artist: 'Pink Floyd', price: 22.99 },
    ];

    albumService.getAll.mockResolvedValue(mockAlbums);

    render(<App />);

    await waitFor(() => {
      expect(screen.queryByText(/loading/i)).not.toBeInTheDocument();
    });

    // Check if albums are displayed
    expect(screen.getByText('The Wall')).toBeInTheDocument();
    expect(screen.getByText('Dark Side of the Moon')).toBeInTheDocument();

    // Check album count
    expect(screen.getByText(/albums \(2\)/i)).toBeInTheDocument();
  });

  test('displays "no albums" message when list is empty', async () => {
    albumService.getAll.mockResolvedValue([]);

    render(<App />);

    await waitFor(() => {
      expect(screen.queryByText(/loading/i)).not.toBeInTheDocument();
    });

    // Check for empty state message
    expect(screen.getByText(/no albums yet/i)).toBeInTheDocument();
  });

  test('displays error message when API fails', async () => {
    albumService.getAll.mockRejectedValue(new Error('API Error'));

    render(<App />);

    await waitFor(() => {
      expect(screen.queryByText(/loading/i)).not.toBeInTheDocument();
    });

    // Check for error message
    expect(screen.getByText(/failed to load albums/i)).toBeInTheDocument();
  });
});
