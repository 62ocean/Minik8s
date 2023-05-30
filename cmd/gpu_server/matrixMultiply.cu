#include <iostream>
#include <stdlib.h>
#include <stdio.h>
#include <cuda_runtime_api.h>

#define M 20
#define N 20
#define THREADS_PER_BLOCK 5

// add function of matrix
__global__ void matmul(int *A, int *B, int *C){
    int row = blockIdx.x*blockDim.x+threadIdx.x;
    int col = blockIdx.y*blockDim.y+threadIdx.y;
    int value = 0;
    for(int k=0;k<N;k++){
        value = value+A[row*N+k]*B[k*N+col];
    }
    C[row*N+col] = value;
}

int main(){
    int *d_A,*d_B,*d_C;
    int size = M*N*sizeof (int);
    cudaMalloc((void **)&d_A,size);
    cudaMalloc((void **)&d_B,size);
    cudaMalloc((void **)&d_C,size);

    // initialize the matrix
    int *A = (int *)malloc(size);
    int *B = (int *)malloc(size);
    int *C = (int *)malloc(size);
    for(int i=0;i<M*N;i++){
        A[i]=i;
        B[i]=i;
    }

    dim3 threadPerBlock(5, 5);
    dim3 blocks(M/threadPerBlock.x,N/threadPerBlock.y);

    // copy data into device
    cudaMemcpy(d_A,A,size,cudaMemcpyHostToDevice);
    cudaMemcpy(d_B,B,size,cudaMemcpyHostToDevice);

    matmul<<<blocks,threadPerBlock>>>(d_A,d_B,d_C);

    // copy back
    cudaMemcpy(C,d_C,size,cudaMemcpyDeviceToHost);

    // print result:
    printf("The matrix add result is:\n");
    for (int i = 0; i < M; i++) {
        for (int j = 0; j < N ; j++) {
            printf("%d ", C[i * N + j]);
        }
        printf("\n");
    }

    cudaFree(d_A);
    cudaFree(d_B);
    cudaFree(d_C);
    free(A);
    free(B);
    free(C);
}