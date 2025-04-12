import Button from './atoms/Button'
import Card from './Card'

function ManageUser() {
  return (
    <Card className="shadow-sm bg-white">
      <div className="flex flex-col">
        <div className="flex justify-between items-center">
          <h2>Manage Users</h2>
          <Button classname="border border-green-400">Add User</Button>
        </div>
      </div>
    </Card>
  )
}

export default ManageUser
