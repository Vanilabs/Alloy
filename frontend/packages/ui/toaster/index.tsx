import { Toast, ToastClose, ToastDescription, ToastProvider, ToastTitle, ToastViewport } from "../toast";
import { X, CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-react';
import { useToast } from "../use-toast";

const alertIcons = {
    success: CheckCircle,
    error: AlertCircle,
    warning: AlertTriangle,
    info: Info,
};

const alertStyles = {
    success: 'bg-green-500/10 border-green-500/30 text-green-400',
    error: 'bg-red-500/10 border-red-500/30 text-red-400',
    warning: 'bg-yellow-500/10 border-yellow-500/30 text-yellow-400',
    info: 'bg-blue-500/10 border-blue-500/30 text-blue-400',
};

const iconStyles = {
    success: 'text-green-400',
    error: 'text-red-400',
    warning: 'text-yellow-400',
    info: 'text-blue-400',
};

export function Toaster() {
    const { toasts } = useToast();

    return (
        <ToastProvider swipeDirection="left">
            {toasts.map(function ({ id, title, description, status, action, ...props }) {

                const Icon = status ? alertIcons[status] : alertIcons.info;

                return (
                    <Toast key={id} {...props} className={`mb-4 max-w-sm border z-30 ${alertStyles[status!!] || alertStyles.info}`}>
                        <div className="flex gap-1">
                            {Icon && <Icon className={`w-5 h-5 mb-1 ${iconStyles[status!!] || iconStyles.info}`} />}
                            <div className="flex flex-col">
                                {title && <ToastTitle>{title}</ToastTitle>}
                                {description && <ToastDescription>{description}</ToastDescription>}
                            </div>
                        </div>
                        {action}
                        <ToastClose />
                    </Toast>
                );
            })}
            <ToastViewport />
        </ToastProvider>
    );
}
