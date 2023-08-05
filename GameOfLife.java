public class GameOfLife {
    private int rows;
    private int cols;
    private boolean[][] grid;

    public GameOfLife(int rows, int cols) {
        this.rows = rows;
        this.cols = cols;
        grid = new boolean[rows][cols];
    }

    // Method to set the initial state of the grid
    public void setInitialState(int[][] initialState) {
        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                grid[i][j] = initialState[i][j] == 1;
            }
        }
    }

    // Method to print the current state of the grid
    public void printGrid() {
        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                System.out.print(grid[i][j] ? "X " : ". ");
            }
            System.out.println();
        }
        System.out.println();
    }

    // Method to calculate the next generation of the grid
    public void nextGeneration() {
        boolean[][] newGrid = new boolean[rows][cols];
        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                int liveNeighbors = countLiveNeighbors(i, j);
                if (grid[i][j]) {
                    // Any live cell with fewer than two live neighbors dies (underpopulation)
                    // Any live cell with two or three live neighbors lives on to the next generation
                    // Any live cell with more than three live neighbors dies (overpopulation)
                    newGrid[i][j] = liveNeighbors == 2 || liveNeighbors == 3;
                } else {
                    // Any dead cell with exactly three live neighbors becomes a live cell (reproduction)
                    newGrid[i][j] = liveNeighbors == 3;
                }
            }
        }
        grid = newGrid;
    }

    // Method to count the number of live neighbors around a cell
    private int countLiveNeighbors(int row, int col) {
        int count = 0;
        for (int i = -1; i <= 1; i++) {
            for (int j = -1; j <= 1; j++) {
                if (i == 0 && j == 0) {
                    continue; // Skip the center cell (the cell itself)
                }
                int newRow = row + i;
                int newCol = col + j;
                if (newRow >= 0 && newRow < rows && newCol >= 0 && newCol < cols) {
                    if (grid[newRow][newCol]) {
                        count++;
                    }
                }
            }
        }
        return count;
    }

    public static void main(String[] args) {
        int rows = 8;
        int cols = 10;
        int[][] initialState = {
            {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
            {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
            {0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
            {0, 0, 0, 1, 1, 1, 0, 0, 0, 0},
            {0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
            {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
            {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
            {0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
        };

        GameOfLife game = new GameOfLife(rows, cols);
        game.setInitialState(initialState);

        for (int generation = 1; generation <= 5; generation++) {
            System.out.println("Generation " + generation + ":");
            game.printGrid();
            game.nextGeneration();
        }
    }
}
