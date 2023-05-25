#include <iostream>
#include <stdlib.h>

#define M 20
#define N 20

// add function of matrix
__global__ void matrix_add(int **A, int **B, int **C){

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
    for(int i=0;i<size;i++){
        A[i]=i;
        B[i]=i;
    }

    // copy data into device
    cudaMemcpy(d_A,a,size,cudaMemcpyHostToDevice)


}