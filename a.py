def multiply(a, b):
      return [
          [a[0][0] * b[0][0] + a[0][1] * b[1][0], a[0][0] * b[0][1] + a[0][1] * b[1][1]],
          [a[1][0] * b[0][0] + a[1][1] * b[1][0], a[1][0] * b[0][1] + a[1][1] * b[1][1]]
      ]
  
def matrix_power(matrix, n):
      if n == 0:
          return [[1, 0], [0, 1]]
      if n == 1:
          return matrix
      
      half = matrix_power(matrix, n >> 1)
      squared = multiply(half, half)
      
      if n & 1:
          return multiply(squared, matrix)
      return squared
  
def fib(n):
      if n <= 1:
          return n
      base_matrix = [[1, 1], [1, 0]]
      result_matrix = matrix_power(base_matrix, n - 1)
      return result_matrix[0][0]
  
n = int(input())
print(fib(n))