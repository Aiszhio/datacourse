import Swal from 'sweetalert2';
import 'sweetalert2/dist/sweetalert2.min.css';

export const showSuccess = (message) => {
    Swal.fire({
        title: 'Успех!',
        text: message,
        icon: 'success',
        confirmButtonText: 'OK',
        toast: true,
        position: 'top-end',
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
    });
};

export const showError = (message) => {
    Swal.fire({
        title: 'Ошибка!',
        text: message,
        icon: 'error',
        confirmButtonText: 'OK',
        toast: true,
        position: 'top-end',
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: true,
    });
};
