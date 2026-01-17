interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: 'primary' | 'secondary'
}

export const Button = ({ variant = 'primary', ...props }: ButtonProps) => {
    const base = 'px-4 py-2 rounded-md font-medium'
    const color = variant === 'primary' ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-800'
    return <button className={`${base} ${color}`} {...props} />
}
