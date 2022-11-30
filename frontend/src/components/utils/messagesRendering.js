export const isSenderTheLoggedInUser = (m, loggedInUserId) => {
    return m.sender[0]._id === loggedInUserId
}