#include <iostream>
#include <stdlib.h>
#include <stdio.h>
#include <cuda_runtime_api.h>

#define M 20
#define N 20
#define THREADS_PER_BLOCK 5

// add function of matrix
__global__ void matrix_add(int *A, int *B, int *C){
    C[threadIdx.x+blockIdx.x*THREADS_PER_BLOCK]=A[threadIdx.x+blockIdx.x*5]+B[threadIdx.x+blockIdx.x*5];
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

    // copy data into device
    cudaMemcpy(d_A,A,size,cudaMemcpyHostToDevice);
    cudaMemcpy(d_B,B,size,cudaMemcpyHostToDevice);

    matrix_add<<<size/THREADS_PER_BLOCK,THREADS_PER_BLOCK>>>(d_A,d_B,d_C);

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
}