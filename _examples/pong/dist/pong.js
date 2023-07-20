// DOM elements
const playingField = document.getElementById('playingField');
const userPaddle = document.getElementById('userPaddle');
const computerPaddle = document.getElementById('computerPaddle');
const ball = document.getElementById('ball');
const userScore = document.getElementById('userScore');
const computerScore = document.getElementById('computerScore');
const message = document.getElementById('message');

// Game variables
let userPaddleY = 0;
let computerPaddleY = 0;
let ballX = 0;
let ballY = 0;
let ballSpeedX = 0;
let ballSpeedY = 0;
let userScoreCount = 0;
let computerScoreCount = 0;
let gameStarted = false;

// Start the game
function startGame() {
  // Initialize game variables
  userPaddleY = playingField.offsetHeight / 2 - userPaddle.offsetHeight / 2;
  computerPaddleY = playingField.offsetHeight / 2 - computerPaddle.offsetHeight / 2;
  ballX = playingField.offsetWidth / 2 - ball.offsetWidth / 2;
  ballY = playingField.offsetHeight / 2 - ball.offsetHeight / 2;
  ballSpeedX = 3;
  ballSpeedY = 3;
  userScoreCount = 0;
  computerScoreCount = 0;
  gameStarted = true;
  
  // Reset paddle and ball positions
  userPaddle.style.top = userPaddleY + 'px';
  computerPaddle.style.top = computerPaddleY + 'px';
  ball.style.left = ballX + 'px';
  ball.style.top = ballY + 'px';
  
  // Reset scores and message
  userScore.textContent = userScoreCount;
  computerScore.textContent = computerScoreCount;
  message.textContent = '';
  
  // Start game loop
  requestAnimationFrame(moveBall);
}

// Move the user's paddle
function moveUserPaddle(event) {
  if (gameStarted) {
    switch (event.key) {
      case 'ArrowUp':
        userPaddleY -= 10;
        break;
      case 'ArrowDown':
        userPaddleY += 10;
        break;
    }
    // Prevent paddle from going out of bounds
    if (userPaddleY < 0) {
      userPaddleY = 0;
    }
    if (userPaddleY > playingField.offsetHeight - userPaddle.offsetHeight) {
      userPaddleY = playingField.offsetHeight - userPaddle.offsetHeight;
    }
    userPaddle.style.top = userPaddleY + 'px';
  }
}

// Update the computer's paddle position to follow the ball
function updateComputerPaddle() {
  if (gameStarted) {
    const paddleCenter = computerPaddleY + computerPaddle.offsetHeight / 2;
    if (paddleCenter < ballY) {
      computerPaddleY += 3;
    } else {
      computerPaddleY -= 3;
    }
    // Prevent paddle from going out of bounds
    if (computerPaddleY < 0) {
      computerPaddleY = 0;
    }
    if (computerPaddleY > playingField.offsetHeight - computerPaddle.offsetHeight) {
      computerPaddleY = playingField.offsetHeight - computerPaddle.offsetHeight;
    }
    computerPaddle.style.top = computerPaddleY + 'px';
  }
}

// Move the ball
function moveBall() {
  if (gameStarted) {
    // Update ball position
    ballX += ballSpeedX;
    ballY += ballSpeedY;
    ball.style.left = ballX + 'px';
    ball.style.top = ballY + 'px';
    
    // Check collision with paddles
    checkCollision();
    
    // Check collision with walls
    if (ballY < 0 || ballY > playingField.offsetHeight - ball.offsetHeight) {
      ballSpeedY *= -1;
    }
    
    // Check if ball is out of bounds
    if (ballX < 0) {
      // Computer scores
      computerScoreCount++;
      computerScore.textContent = computerScoreCount;
      resetBall();
    }
    if (ballX > playingField.offsetWidth - ball.offsetWidth) {
      // User scores
      userScoreCount++;
      userScore.textContent = userScoreCount;
      resetBall();
    }
    
    // Check win condition
    checkWin();
    
    // Continue game loop
    requestAnimationFrame(moveBall);
  }
}

// Check collision with paddles
function checkCollision() {
  if (ballX < userPaddle.offsetWidth && ballY + ball.offsetHeight > userPaddleY && ballY < userPaddleY + userPaddle.offsetHeight) {
    // Ball collides with user's paddle
    ballSpeedX *= -1;
  }
  if (ballX + ball.offsetWidth > playingField.offsetWidth - computerPaddle.offsetWidth && ballY + ball.offsetHeight > computerPaddleY && ballY < computerPaddleY + computerPaddle.offsetHeight) {
    // Ball collides with computer's paddle
    ballSpeedX *= -1;
  }
}

// Reset ball position
function resetBall() {
  ballX = playingField.offsetWidth / 2 - ball.offsetWidth / 2;
  ballY = playingField.offsetHeight / 2 - ball.offsetHeight / 2;
  ballSpeedX *= -1;
  ballSpeedY = Math.random() > 0.5 ? 3 : -3;
  ball.style.left = ballX + 'px';
  ball.style.top = ballY + 'px';
}

// Check win condition
function checkWin() {
  if (userScoreCount === 10) {
    // User wins
    gameStarted = false;
    message.textContent = 'You win!';
  }
  if (computerScoreCount === 10) {
    // Computer wins
    gameStarted = false;
    message.textContent = 'Computer wins!';
  }
}

// Event listeners
document.addEventListener('keydown', moveUserPaddle);

// Start the game
startGame();
