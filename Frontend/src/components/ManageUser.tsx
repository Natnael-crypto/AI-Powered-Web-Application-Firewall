import Button from './atoms/Button'
import Card from './Card'
interface manageUserProps {
  toggleAddUser: () => void
}
function ManageUser({toggleAddUser}: manageUserProps) {
  return (
    <Card className="shadow-sm bg-white">
      <div className="flex flex-col">
        <div className="flex justify-between items-center">
          <h2>Manage Users</h2>
          <Button classname="border border-gray-600 text-white" onClick={toggleAddUser}>
            Add User
          </Button>
        </div>
      </div>
    </Card>
  )
}

export default ManageUser
